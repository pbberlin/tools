// Package distrib enables distributed processing of
// any slice of structs implementing the Worker
// interface.
//
// Distrib contains the simplest pipeline with three stages.
// Stage 2 is fanned out, executing Work().
//
// All send and receive operations are clad in
// select statements with a <-fin branch,
// ensuring clear exit of all goroutines.
//
// Distrib processes packages of type Packet.
// The actual work is always contained in Packet.Work().
//
// Distrib has two exit modes.
// CollectRemainder or inversely "FlushAndAbandon".
//
// FlushAndAbandon was supposed to set inn and out to nil,
// blocking every further communication.
// It then closes fin, causing all goroutines to exit.
// Setting inn and out to nil alarmed the race detector.
// Instead we signal stage1 to stop with "lcnt = -10".
//
// CollectRemainder equally prevents further feeding of the pipeline.
// It then waits until all sent packets
// have been received at stage3,
// relying on the synchronized sent counter.
//
// Distrib provides a timeout for Packet.Work().
// Timing out of Work() requires spawning a goroutine
// and an individual receiver channel for each packet.
// Upon timeout the channel receiving the work is
// closed.
// The timed out workers will hang around in blocked mode
// until they are flushed with close(fin) at exit.
//
// Stage3 of Distrib is a *synchroneous* for-select loop,
// redeeming us from waitgroup.Wait...
//
// Signalling stage1 to stage3 that packets are exhausted
// happens by setting Options.Want to zero.
// Options.Want is therefore not usable as cutoff
// for the returns slice.
//
// Signalling stage3 to stage1 that loading should be
// stopped, happens by setting lcnt to -10.
// Both need to be synchronized.
//
// Todo/Caveat: CollectRemainder==false && len(jobs) << Want
// leads to premature flushing
// Better use CollectRemainder==true
//
// We do not use loghttp.FuncBufUniv,
// since it is not threadseave.
// Instead we us log.Printf to write into a byte buffer,
// which is returned to caller.
package distrib

import (
	"bytes"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/pbberlin/tools/net/http/loghttp"
)

// Worker interface is a narrow interface
// executing the work, that is to be distributed.
type Worker interface {
	Work()
}

// Packet is a struct, that is passed around between stages.
// Packet.Work() is evoked during stage2.
type Packet struct {
	TaskID   int
	WorkerID int
	Worker   // Anonymous interface; any struct that has a Work() method
}

func (w Packet) String() string {
	return fmt.Sprintf(" packet#%-2v  Wkr%v", w.TaskID, w.WorkerID)
}

type Options struct {
	NumWorkers       int           //
	Want             int32         // Maximum results before pipeline is torn down
	TimeOutDur       time.Duration // Max duration of Work() before timeout + abandonment.
	CollectRemainder bool          // Upon exit: Wait for remaining packets to reach stage3, or flush them where they are.
	TailingSleep     bool          // Upon exit: Wait a short while, checking the return of all goroutines.

	Logger loghttp.FuncBufUniv // not thread safe, but for debugging appeninge
}

var defaultOptions = Options{
	NumWorkers:       6,
	Want:             int32(10),
	TimeOutDur:       time.Millisecond * 14,
	CollectRemainder: true,
	TailingSleep:     false,
}

func NewDefaultOptions() Options {
	return defaultOptions
}

// Distrib builds the pipleline and processes the packets.
func Distrib(jobs []Worker, opt Options) ([]*Packet, *bytes.Buffer) {

	var b = new(bytes.Buffer)
	// var lnp = log.New(os.Stderr, "", 0) // logger no prefix; os.Stderr shows up in appengine devserver; os.Stdout does not

	var lnp = log.New(b, "", 0)
	var lpf = lnp.Printf // shortcut
	// var lpf = opt.Logger

	inn := make(chan *Packet) // stage1 => stage2
	out := make(chan *Packet) // stage2 => stage3

	fin := make(chan struct{}) // flush all

	lcnt := int32(0) // load counter; always incrementing; except: to signal stop loading from stage3 to stage1
	sent := int32(0) // sent packages - might get decremented for timed out packets
	recv := int32(0) // received packages - converges against sent packages - unless in flushing mode

	ticker := time.NewTicker(10 * time.Millisecond)
	tick := ticker.C // channel from ticker

	var returns = make([]*Packet, 0, int(opt.Want)+5)

	if len(jobs) < 1 {
		lpf("empty jobs slice")
		return returns, b
	}

	packets := make([]*Packet, 0, len(jobs))
	for i, job := range jobs {
		pack := &Packet{}
		pack.Worker = job
		pack.TaskID = i
		packets = append(packets, pack)
	}

	if opt.NumWorkers > len(packets) && len(packets) > 0 {
		opt.NumWorkers = len(packets)
		lpf("num workers curtailed to |%v|", opt.NumWorkers)
	}
	lpf("want |%v| with num workers |%v| from |%v|%v| jobs|packets",
		opt.Want, opt.NumWorkers, len(jobs), len(packets))

	//
	// stage 1
	go func() {
		for {

			idx := int(atomic.LoadInt32(&lcnt))

			if idx > len(packets)-1 { // signal to stage3
				lpf("=== input packets exhausted at %v ===", idx)
				atomic.StoreInt32(&opt.Want, int32(0))
				return
			}
			if idx < 0 { // signal from stage3
				lpf("=== loading stage 1 terminated ===")
				return
			}
			select {
			case inn <- packets[idx]:
				atomic.AddInt32(&lcnt, 1)
				atomic.AddInt32(&sent, 1) // sent++, sent can be decremented later on, therefore distinct var "lcnt"
			case <-fin:
				return
			}
		}
	}()

	//
	// stage 2
	for i := 0; i < opt.NumWorkers; i++ {
		go func(wkrID int) {
		NextPacket:
			for {
				timeout := time.After(opt.TimeOutDur)
				select {
				case packet := <-inn:

					// Wrap work into a go routine.
					// It puts the result into chan res.
					res := make(chan *Packet)
					workWrap := func(pck *Packet) { // packet given as argument to prevent race cond race_1_b
						pck.WorkerID = wkrID // signature
						pck.Work()
						select {
						case res <- pck:
						case <-fin:
						}
					}
					go workWrap(packet)

					//
					// Now put workWrap() in competition with timeout.
					select {
					case completed := <-res:
						select {
						case out <- completed:
						case <-fin:
							return
						}
					case <-timeout:
						atomic.AddInt32(&sent, -1)
						lpf("TOUT  snt%2v  %v", atomic.LoadInt32(&sent), *packet) // race_1_b
						// => stage 3 has to check recv >= sent in separate select-tick branch
						//    because no WLoad packet is sent-received upon this timeout.
						break NextPacket // skipping this packet
					}

				case <-fin:
					return
				}
			}
		}(i)
	}

	//
	// stage 3
	// synchroneous
	func() {
		for {
			select {

			// case <-time.After(22 * time.Second):
			// 	lpf("=== TOTAL TIMEOUT  ===")
			// 	func() {
			// 		defer func() {
			// 			recover()
			// 		}()
			// 		close(fin)
			// 	}()
			// 	return

			case <-tick:
				// Exit after collecting remainder
				// Sent might be decremented in timed out works
				if opt.CollectRemainder &&
					recv >= atomic.LoadInt32(&opt.Want) &&
					recv >= atomic.LoadInt32(&sent) &&
					true {
					lpf("=== enough on tick ===")
					return
				}

			case packet := <-out:

				recv++

				lpf("rcv%-2v of %-2v snt%-2v  %v  %v", recv, atomic.LoadInt32(&opt.Want), atomic.LoadInt32(&sent), packet, time.Now().Format("05.000"))
				// lpf("  %v", stringspb.IndentedDump(packet.Workload))

				returns = append(returns, packet)

				if recv >= atomic.LoadInt32(&opt.Want) {

					if recv == atomic.LoadInt32(&opt.Want) {
						lpf("=== enough on receive ===")
					}

					// inn = nil  // race detector objected to this line
					atomic.StoreInt32(&lcnt, -10) // signalling stop loading to stage1

					// Exit immediately
					if !opt.CollectRemainder {
						lpf("=== flush all stages and abandon remaining results ===")
						// out = nil // race detector objected to this line
						close(fin)
						if opt.TailingSleep {
							time.Sleep(60 * time.Millisecond) // few messages from exiting goroutines might occur
						}
						return
					}

				}

			case <-fin:
				return
			}
		}
	}()

	lpf("=== cleanup ===")

	ticker.Stop() // cleaning up the ticker really necessary?

	if opt.CollectRemainder {
		close(fin) // flush remaining packets whereever they are
		if opt.TailingSleep {
			time.Sleep(60 * time.Millisecond) // 1 or 2 messages from timed out workWrap might occur
		}
	}

	return returns, b

}
