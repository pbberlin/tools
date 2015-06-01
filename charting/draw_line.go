package charting

import (
	"image"
	"image/color"	
	

	"math"
	"log"

	"github.com/pbberlin/tools/util"
	
)


// Blend a new pixel over an old pixel - heeding their alpha chan values
// 
// algorithm NOT according to http://en.wikipedia.org/wiki/Alpha_compositing
//   but by my own trial and error
func blendPixelOverPixel(ic_old,ic_new uint8, al_new float64)(c_res uint8) {

	al_old := float64(1); _=al_old
	c_old  := float64(ic_old)
	c_new  := float64(ic_new)

	algo1 := c_old*(1-al_new)   +   c_new*al_new
	c_res =  uint8( util.Min( util.Round(algo1),255) )
	//log.Printf("\t\t %3.1f +  %3.1f  = %3.1f", c_old*(1-al_new),c_new*al_new, algo1)

	return 
}



func funcSetPixler(col color.RGBA, img *image.RGBA )( func( addr int,dist float64) ){

	// 	4*400*300
	r  := img.Rect
	p0 := r.Min
	p1 := r.Max
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	maxPix := dx * dy * 4
	log.Printf("\tfuncSetPixler  BxH: %vx%v  Size:%v (%v)",dx,dy,maxPix,len(img.Pix) )	
	
	return func(addr int,dist float64){
		
		//log.Printf("\t%v<%v",addr,maxPix )	
		if addr > (maxPix-4)  ||  addr < 0 {
			log.Printf("\t%v<%v !  OVERFLOW! ",addr,maxPix )	
			return
		}

		// dist ranges from 0 to 1.5
		if dist < 0.0 { 
			dist = 0
		}
		
		sharpness := 0.9  // < 1 => more blurred ; otherwise more pixely
		ba :=  math.Pow( 1 - (dist * 2/3), sharpness )

		//log.Printf("\tbef: %3d %3d %3d",img.Pix[addr+0],img.Pix[addr+1],img.Pix[addr+2])
		//log.Printf("\tcol: %3d %3d %3d | %1.3f => %1.3f",col.R,col.G,col.B,dist,ba)

		img.Pix[addr+0] = blendPixelOverPixel(img.Pix[addr+0],col.R,  ba) 
		img.Pix[addr+1] = blendPixelOverPixel(img.Pix[addr+1],col.G,  ba) 
		img.Pix[addr+2] = blendPixelOverPixel(img.Pix[addr+2],col.B,  ba) 

		//log.Printf("\taft: %3d %3d %3d\n\n",img.Pix[addr+0],img.Pix[addr+1],img.Pix[addr+2])

	}
	
}







// https://courses.engr.illinois.edu/ece390/archive/archive-f2000/mp/mp4/anti.html
func FuncDrawLiner(lCol color.RGBA, img *image.RGBA )( func( P_next image.Point, lCol color.RGBA, img *image.RGBA )  ){

	var P_last image.Point = image.Point{-1111,-1111}

	r  := img.Rect
	p0 := r.Min
	p1 := r.Max
	imgWidth := p1.X - p0.X


	return func (P_next image.Point, lCol color.RGBA, img *image.RGBA ){

		var P0, P1 image.Point

		if P_last.X == -1111  &&  P_last.Y == -1111{
			P_last = P_next
			return	
		} else {
			P0 = P_last	
			P1 = P_next
			P_last = P_next	
		}
		
		
		log.Printf("draw_line_start---------------------------------")
	
		x0, y0 := P0.X, P0.Y
		x1, y1 := P1.X, P1.Y
	
		
		bpp := 4  // bytes per pixel
	
		addr := (y0*imgWidth+x0)*bpp
		dx   := x1-x0
		dy   := y1-y0
	
	
		var du, dv,u ,v int
		var uincr int = bpp
		var vincr int = imgWidth*bpp
	
		
		// switching to (u,v) to combine all eight octants
		if  util.Abs(dx) > util.Abs(dy) {
			du = util.Abs(dx)
			dv = util.Abs(dy)
			u = x1
			v = y1
			uincr = bpp
			vincr = imgWidth*bpp
			if dx < 0 {uincr = -uincr}
			if dy < 0 {vincr = -vincr}
		} else {
			du = util.Abs(dy)
			dv = util.Abs(dx)
			u = y1
			v = x1
			uincr = imgWidth*bpp
			vincr = bpp
			if dy < 0 {uincr = -uincr}
			if dx < 0 {vincr = -vincr}
		}
		log.Printf("draw_line\tu %v - v %v - du %v - dv %v - uinc %v - vinc %v ", u, v, du, dv, uincr, vincr)
		
		// uend	  :=  u + 2 * du
		// d	     := (2 * dv) - du		// Initial value as in Bresenham's 
		// incrS   :=  2 *  dv				// Δd for straight increments 
		// incrD   :=  2 * (dv - du)	   // Δd for diagonal increments 
		// twovdu  :=  0						// Numerator of distance starts at 0 
	
	
		// I have NO idea why - but unless I use -1- 
		//   instead of the orginal -2- as factor,
		//   all lines are drawn DOUBLE the intended size
		//   THIS is how it works for me:
		uend	  :=  u + 1 * du
		d	     := (1 * dv) - du		// Initial value as in Bresenham's 
		incrS   :=  1 *  dv				// Δd for straight increments 
		incrD   :=  1 * (dv - du)	   // Δd for diagonal increments 
		twovdu  :=  0						// Numerator of distance starts at 0 
	
						
	
		log.Printf("draw_line\tuend %v - d %v - incrS %v - incrD %v - twovdu %v", uend, d, incrS, incrD, twovdu)
	
	
		tmp     := float64(du*du + dv*dv)
		invD	  := 1.0 / (2.0*math.Sqrt( tmp ))   /* Precomputed inverse denominator */
		invD2du := 2.0 * (  float64(du)*invD)	   /* Precomputed constant */
	
		log.Printf("draw_line\tinvD %v - invD2du %v", invD, invD2du)
	
		cntr := -1
		
		setPix := funcSetPixler(lCol,img)
		
		for{
			cntr++
			//log.Printf("==lp%v ", cntr )
	
			// Ensure that addr is valid
			ftwovdu:= float64(twovdu)
			setPix(addr - vincr, invD2du + ftwovdu*invD)
			setPix(addr		  ,	        ftwovdu*invD)
			setPix(addr + vincr, invD2du - ftwovdu*invD)
		
	
			if (d < 0){
				/* choose straight (u direction) */
				twovdu = d + du
				d = d + incrS
			} 	else 	{
				/* choose diagonal (u+v direction) */
				twovdu = d - du
				d = d + incrD
				v = v+1
				addr = addr + vincr
			}
			u = u+1
			addr = addr+uincr
			
			if u > uend {break}
		} 
	
		log.Printf("draw_line_end---------------------------------")
	
	}
}	





