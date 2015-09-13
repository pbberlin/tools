package distrib

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

var lnp = log.New(os.Stdout, "", 0) // logger no prefix
var lpf = lnp.Printf

type Worker interface {
	Work()
}

type WLoad struct {
	TaskID   int
	WorkerID int
	Workload Worker
	// Result   int
	// Func     func(WLoad) WLoad
}

type Options struct {
	CollectRemainder bool
	NumWorkers       int
	Want             int32
	TimeOutDur       time.Duration
	TailingSleep     bool // Check for returning goroutines
}

var DefaultOptions = Options{
	CollectRemainder: true,
	NumWorkers:       6,
	Want:             int32(10),
	TimeOutDur:       time.Millisecond * 15,
	TailingSleep:     false,
}

func (w WLoad) String() string {
	return fmt.Sprintf(" packet#%-2v  Wkr%v", w.TaskID, w.WorkerID)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Distrib contains the simplest pipeline of three stages.
//
// All send and receive operations are clad in
// select statements with a <-fin branch,
// ensuring clear exit of all goroutines.
//
// Distrib processes packages of type WLoad.
// The actual work is always contained in WLoad.Func().
// WLoad.Func() always takes a packet and returns a packet.
//
// Distrib has two exit modes.
// CollectRemainder or FlushAndAbandon.
//
// FlushAndAbandon sets inn and out to nil,
// blocking every further communication.
// It then closes fin, causing all goroutines to exit.
//
// CollectRemainder only block inn,
// preventing further feeding of the pipeline.
// It then waits until all sent packets
// have been received at stage3,
// relying on the synchronized sent counter.
//
// Distrib provides a timeout for WLoad.Func().
// Timeout for workWrap requires spawning of a goroutine
// and an indivudual receiver channel for each packet.
// Upon timeout the channel receiving the work is
// closed.
// The timed out workers will hang around in blocked mode
// until they are flushed with close(fin) at the end.
//
// Stage 3 of Distrib is a synchroneous loop,
// delivering us the waitgroup.Wait...
//
// Signalling stage 1 to stage 3 that packets are exhausted
// happens by setting want to zero.
//
// Signalling stage 3 to stage 1 that loading should be
// stopped, happens by setting lcnt to -10.
// Both need to be synchronized.
//
//
// Todo: Change all channels and the slices to *pointers* of WLoad, reducing memory copy
func Distrib(packets []WLoad, opt Options) []WLoad {

	inn := make(chan WLoad) // stage1 => stage2
	out := make(chan WLoad) // stage2 => stage3

	fin := make(chan struct{}) // flush all

	lcnt := int32(0) // load counter; always incrementing; except: to signal stop loading from downstream
	sent := int32(0) // sent packages - might get decremented for timed out work
	recv := int32(0) // received packages - converges against sent packages - unless in flushing mode

	ticker := time.NewTicker(10 * time.Millisecond)
	tick := ticker.C // channel from ticker

	var ret = make([]WLoad, 0, int(opt.Want)+5)

	for i, _ := range packets {
		if packets[i].TaskID == 0 {
			packets[i].TaskID = i
		}
	}

	//
	// stage 1
	go func() {
		for {

			idx := int(atomic.LoadInt32(&lcnt))

			if idx > len(packets)-1 {
				lpf("=== input packets exhausted ===")
				atomic.StoreInt32(&opt.Want, int32(0))
				return
			}
			if idx < 0 { // signal from downstream
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
			for {
				timeout := time.After(opt.TimeOutDur)
			MarkX:
				select {
				case packet := <-inn:

					// Wrap work into a go routine.
					// It puts the result into chan res.
					res := make(chan WLoad)
					workWrap := func(pck WLoad) { // packet as argument to prevent race cond race_1_b
						pck.WorkerID = wkrID // signature
						pck.Workload.Work()
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
						lpf("TOUT  snt%2v  %v", atomic.LoadInt32(&sent), packet) // race_1_b
						// => stage 3 has to check recv >= sent in separate select-tick branch
						//    because no WLoad packet is sent-received upon this timeout.
						break MarkX // skipping this packet
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

			case <-tick:
				// Exit after collecting remainder
				// Sent might be decremented in timed out works
				if opt.CollectRemainder &&
					recv >= atomic.LoadInt32(&opt.Want) &&
					recv >= atomic.LoadInt32(&sent) &&
					true {
					return
				}

			case packet := <-out:

				recv++

				lpf("rcv%-2v snt%-2v  %v", recv, atomic.LoadInt32(&sent), packet)
				// lpf("  %v", stringspb.IndentedDump(packet.Workload))
				if recv <= atomic.LoadInt32(&opt.Want) {
					ret = append(ret, packet)
				}

				if recv >= atomic.LoadInt32(&opt.Want) {

					if recv == atomic.LoadInt32(&opt.Want) {
						lpf("=== enough ===")
					}

					// inn = nil  // race detector objected to this line
					atomic.StoreInt32(&lcnt, -10) // signalling stop loading to stage 1

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

	ticker.Stop() // cleaning up really neccessary?

	if opt.CollectRemainder {
		close(fin) // flush
		if opt.TailingSleep {
			time.Sleep(60 * time.Millisecond) // 1 or 2 messages from timed out workWrap might occur
		}
	}

	return ret

}
