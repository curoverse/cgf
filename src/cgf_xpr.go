package main

import "fmt"
import "./dlug"

//import "crypto/md5"

type cgfintermediate struct {
  step  [][]int
  seq   [][]string
  varid [][]int
  span  [][]int
  loq   [][]bool
  nocall_start_len [][][]int

  tilemap_key string
  tilemap_pos int
  loq_flag bool
  tot_span int
}

type offsetintermediate struct {
  recno int
  stride int
  tilepos []int
  tilemap []int

  offset_idx []int
  tilepos_idx []int

  final_overflow_flag []bool
  span_flag []bool
}

type finaloffsetintermediate struct {
  recno int
  tilepos []int
  variant_ints []int
}

type vectorelement struct {
  canon_flag bool
  cache_flag bool
  ovf_cache_flag bool
  ovf_flag bool
  fin_ovf_flag bool
  span_flag bool

  loq_flag bool

  knot cgfintermediate
  hexit_pos int
  vec_pos int
}

type loqintermediate struct {
  loqinfo_bytecount int
  count         int
  code          int
  stride        int
  tilepos       []int

  homflag       []bool
  loqinfo_ints  []int

  loq_flag      []bool
}


/*
type loqintermediate struct {
  loq_bv []byte
  length int
  count int
  code int
  stride int
  offset_idx []int
  tilepos_idx []int
  loq_hom_flag []byte
}
*/

//========================================================================
//========================================================================
//========================================================================
//========================================================================
//========================================================================

func vectorelement_ovf_count(prep_vector []vectorelement, st, n int) int {
  ovf_count := 0
  for i:=st; i<(st+n); i++ {
    if prep_vector[i].ovf_flag || prep_vector[i].ovf_cache_flag {
      ovf_count++
    }
  }
  return ovf_count
}


func _knot_tot_span(knot *cgfintermediate) int {
  sp := [2]int{}
  for allele:=0; allele<2; allele++ {
    for i:=0; i<len(knot.span[allele]); i++ {
      sp[allele] += knot.span[allele][i]
    }
  }

  max_span := sp[0]
  if sp[1]>sp[0] { max_span = sp[1] }

  knot.tot_span = max_span
  return max_span
}

// return span
//
func _add_knot(knot *cgfintermediate, allele, step_idx int, ti TileInfo, sglf *SGLF) (int,error) {
  if len(ti.NocallStartLen)>0 {
    knot.loq[allele] = append(knot.loq[allele], true)
    knot.loq_flag = true
  } else {
    knot.loq[allele] = append(knot.loq[allele], false)
  }

  sglf_info := SGLFInfo{}
  var ok bool

  // sglf_info only holds a valid path and step
  //
  if step_idx>0 {
    sglf_info,ok = sglf.PfxTagLookup[ti.PfxTag]
  } else {
    sglf_info,ok = sglf.SfxTagLookup[ti.SfxTag]
  }

  if !ok {
    return -1,fmt.Errorf("could not find prefix (%s) in sglf (allele_idx %d, step_idx %d (%x))\n",
      ti.PfxTag, 0, step_idx, step_idx)

  }

  path := sglf_info.Path
  step := sglf_info.Step

  // We need to search for the variant in the Lib to find
  // the rest of the information, including span
  //
  var_idx,e := lookup_variant_index(ti.Seq, sglf.Lib[path][step])
  if e!=nil { return -1,e }

  sglf_info = sglf.LibInfo[path][step][var_idx]
  span := sglf_info.Span

  seq := sglf.Lib[path][step][var_idx]

  knot_allele_idx := len(knot.varid[allele])

  knot.seq[allele] = append(knot.seq[allele], seq)
  knot.varid[allele] = append(knot.varid[allele], var_idx)
  knot.span[allele] = append(knot.span[allele], span)
  knot.step[allele] = append(knot.step[allele], step)

  nc_vec := make([]int,  0, 1024)
  nc_vec = append(nc_vec, ti.NocallStartLen...)
  dummy := [][]int{}
  dummy = append(dummy, []int{})
  knot.nocall_start_len[allele] = append(knot.nocall_start_len[allele], dummy...)
  knot.nocall_start_len[allele][knot_allele_idx] = append(knot.nocall_start_len[allele][knot_allele_idx], nc_vec...)

  return sglf_info.Span,nil
}

func _init_knot(knot *cgfintermediate) {
  knot.seq = make([][]string, 2)
  knot.varid = make([][]int, 2)
  knot.span = make([][]int, 2)
  knot.step = make([][]int, 2)
  knot.loq = make([][]bool, 2)
  knot.nocall_start_len = make([][][]int, 2)
}

//=====         __  __          _   _       _                               _ _       _
//=====   ___  / _|/ _|___  ___| |_(_)_ __ | |_ ___ _ __ _ __ ___   ___  __| (_) __ _| |_ ___
//=====  / _ \| |_| |_/ __|/ _ \ __| | '_ \| __/ _ \ '__| '_ ` _ \ / _ \/ _` | |/ _` | __/ _ \
//===== | (_) |  _|  _\__ \  __/ |_| | | | | ||  __/ |  | | | | | |  __/ (_| | | (_| | ||  __/
//=====  \___/|_| |_| |___/\___|\__|_|_| |_|\__\___|_|  |_| |_| |_|\___|\__,_|_|\__,_|\__\___|
//=====

func offsetintermediate_cmp(ofsi0, ofsi1 offsetintermediate) error {
  if ofsi0.stride != ofsi1.stride { return fmt.Errorf("stride mismatch") }
  if len(ofsi0.tilepos) != len(ofsi1.tilepos) { return fmt.Errorf("tilepos length mismatch") }
  if len(ofsi0.tilemap) != len(ofsi1.tilemap) { return fmt.Errorf("tilemap length mismatch") }
  if len(ofsi0.final_overflow_flag) != len(ofsi1.final_overflow_flag) { return fmt.Errorf("final_overflow_flag mismatch") }
  if len(ofsi0.span_flag) != len(ofsi1.span_flag) {
    return fmt.Errorf("span_flag mismatch (%v != %v)", len(ofsi0.span_flag), len(ofsi1.span_flag))
  }

  for i:=0; i<len(ofsi0.tilemap); i++ {
    if ofsi0.tilemap[i] != ofsi1.tilemap[i] {
      return fmt.Errorf( fmt.Sprintf("tilemap mismatch at %d, %d != %d", i, ofsi0.tilemap[i], ofsi1.tilemap[i]) )
    }
  }

  for i:=0; i<len(ofsi0.final_overflow_flag); i++ {
    if ofsi0.final_overflow_flag[i] != ofsi1.final_overflow_flag[i] {
      return fmt.Errorf("final_overflow_flag mismatch at %d, %v != %v", i, ofsi0.final_overflow_flag[i], ofsi1.final_overflow_flag[i])
    }
  }

  for i:=0; i<len(ofsi0.span_flag); i++ {
    if ofsi0.span_flag[i] != ofsi1.span_flag[i] {
      return fmt.Errorf("span_flag mismatch at %d, %v != %v", i, ofsi0.span_flag[i], ofsi1.span_flag[i])
    }
  }

  return nil
}

func offsetintermediate_from_bytes(b []byte) offsetintermediate {
  ofsi := offsetintermediate{}

  var dummy uint64
  var dn int

  n:=0

  // Length
  dummy = byte2uint64(b[n:n+8])
  n+=8

  NRec := dummy
  ofsi.recno = int(NRec)

  // Stride
  dummy = byte2uint64(b[n:n+8])
  n+=8

  stride := dummy
  ofsi.stride = int(stride)

  // MapByteCount
  dummy = byte2uint64(b[n:n+8])
  n+=8

  mapbytecount := int(dummy)

  ofs_len := int((NRec+stride-1)/stride)

  ofsi.final_overflow_flag = make([]bool, 0, 1024)
  ofsi.span_flag = make([]bool, 0, 1024)
  ofsi.tilepos_idx = make([]int, 0, 1024)

  ofsi.offset_idx = make([]int, 0, 1024)
  for i:=0; i<ofs_len; i++ {
    dummy = byte2uint64(b[n:n+8])
    n+=8

    ofsi.offset_idx = append(ofsi.offset_idx, int(dummy))
  }

  ofsi.tilepos_idx = make([]int, 0, 1024)
  for i:=0; i<ofs_len; i++ {
    dummy = byte2uint64(b[n:n+8])
    n+=8

    ofsi.tilepos_idx = append(ofsi.tilepos_idx, int(dummy))
  }

  for i:=0; i<int(NRec); i++ {
    ofsi.tilepos = append(ofsi.tilepos, ofsi.tilepos_idx[i/int(stride)])
  }

  read_rec := 0
  map_byte_pos:=0
  for map_byte_pos < mapbytecount {
    dummy,dn = dlug.ConvertUint64(b[n:])
    map_byte_pos += dn
    n+=dn


    is_span := false
    if dummy == 1024 {
      dummy = 0
      is_span = true
    }

    is_fin_ovf := false
    if dummy == 1025 {
      dummy = 0
      is_fin_ovf = true
    }

    ofsi.tilemap = append(ofsi.tilemap, int(dummy))

    ofsi.final_overflow_flag = append(ofsi.final_overflow_flag, is_fin_ovf)
    ofsi.span_flag = append(ofsi.span_flag, is_span)

    read_rec++

  }

  /*
  fmt.Printf("converted tilemap %d (read_rec %d, |span_flag| = %d):\n", len(ofsi.tilemap), read_rec, len(ofsi.span_flag))
  for i:=0; i<len(ofsi.tilemap); i++ {
    fmt.Printf("  [%d] %d\n", i, ofsi.tilemap[i])
  }
  */

  return ofsi
}

func bytes_from_offsetintermediate(ofsi offsetintermediate) []byte {
  buf := make([]byte, 64)
  offset_bytes := make([]byte, 0, 1024)

  // number of records
  tobyte64(buf, uint64(len(ofsi.tilepos)))
  offset_bytes = append(offset_bytes, buf[0:8]...)

  // stride
  tobyte64(buf, uint64(ofsi.stride))
  offset_bytes = append(offset_bytes, buf[0:8]...)

  offset_idx := make([]uint64, 0, 1024)
  tilepos_idx := make([]uint64, 0, 1024)

  // construct map bytes, record offset and tilepos index
  // along the way.
  //
  map_bytes := make([]byte, 0, 1024)
  for i:=0; i<len(ofsi.tilepos); i++ {
    if i%256 == 0 {
      offset_idx = append(offset_idx, uint64(len(map_bytes)))
      tilepos_idx = append(tilepos_idx, uint64(ofsi.tilepos[i]))
    }

    val := ofsi.tilemap[i]
    if ofsi.span_flag[i] { val = 1024 }
    if ofsi.final_overflow_flag[i] { val = 1025 }

    mbytes := dlug.MarshalUint64(uint64(val))
    map_bytes = append(map_bytes, mbytes...)
  }

  // MapByteCount
  tobyte64(buf, uint64(len(map_bytes)))
  offset_bytes = append(offset_bytes, buf[0:8]...)

  // offset
  for i:=0; i<len(offset_idx); i++ {
    tobyte64(buf, offset_idx[i])
    offset_bytes = append(offset_bytes, buf[0:8]...)
  }

  // tilepos
  for i:=0; i<len(tilepos_idx); i++ {
    tobyte64(buf, tilepos_idx[i])
    offset_bytes = append(offset_bytes, buf[0:8]...)
  }

  // map
  offset_bytes = append(offset_bytes, map_bytes...)

  return offset_bytes

}

func construct_offset_intermediate(ctx *CGFContext, prep_vector []vectorelement) offsetintermediate {
  ofsi := offsetintermediate{}

  ofsi.tilepos = make([]int, 0, 1024)
  ofsi.tilemap = make([]int, 0, 1024)

  ofsi.offset_idx = make([]int, 0, 1024)
  ofsi.tilepos_idx = make([]int, 0, 1024)

  ofsi.final_overflow_flag = make([]bool, 0, 1024)
  ofsi.span_flag = make([]bool, 0, 1024)

  for i:=0; i<len(prep_vector); i++ {
    if prep_vector[i].canon_flag { continue }
    if prep_vector[i].ovf_flag {
      ofsi.tilepos = append(ofsi.tilepos, i)
      ofsi.tilemap = append(ofsi.tilemap, prep_vector[i].knot.tilemap_pos)

      tf := prep_vector[i].fin_ovf_flag
      //tf := false
      //if prep_vector[i].knot.tilemap_pos > 1023 { tf = true }
      ofsi.final_overflow_flag = append(ofsi.final_overflow_flag, tf)

      tf = false
      if prep_vector[i].span_flag { tf = true }
      ofsi.span_flag = append(ofsi.span_flag, tf)

    }
  }

  ofsi.stride = 256

  return ofsi
}


//========================================================================
//========================================================================
//========================================================================
//========================================================================
//========================================================================

//====    __ _             _        __  __          _   _       _                               _ _       _
//====   / _(_)_ __   __ _| | ___  / _|/ _|___  ___| |_(_)_ __ | |_ ___ _ __ _ __ ___   ___  __| (_) __ _| |_ ___
//====  | |_| | '_ \ / _` | |/ _ \| |_| |_/ __|/ _ \ __| | '_ \| __/ _ \ '__| '_ ` _ \ / _ \/ _` | |/ _` | __/ _ \
//====  |  _| | | | | (_| | | (_) |  _|  _\__ \  __/ |_| | | | | ||  __/ |  | | | | | |  __/ (_| | | (_| | ||  __/
//====  |_| |_|_| |_|\__,_|_|\___/|_| |_| |___/\___|\__|_|_| |_|\__\___|_|  |_| |_| |_|\___|\__,_|_|\__,_|\__\___|

func finaloffsetintermediate_cmp(fofsi0, fofsi1 finaloffsetintermediate) error {
  if len(fofsi0.tilepos)!=len(fofsi1.tilepos) {
    return fmt.Errorf( fmt.Sprintf("tilepos length mismatch: %d != %d", len(fofsi0.tilepos), len(fofsi1.tilepos)) )
  }

  if len(fofsi0.variant_ints) != len(fofsi1.variant_ints) {
    return fmt.Errorf( fmt.Sprintf("variant_ints length mismatch: %d != %d", len(fofsi0.variant_ints), len(fofsi1.variant_ints)) )
  }

  for i:=0; i<len(fofsi0.variant_ints); i++ {
    if fofsi0.variant_ints[i] != fofsi1.variant_ints[i] {
      return fmt.Errorf( fmt.Sprintf("variant_ints mismatch at %d: %d != %d", i, fofsi0.variant_ints[i], fofsi1.variant_ints[i]) )
    }
  }

  return nil
}

func finaloffsetintermediate_from_bytes(b []byte) finaloffsetintermediate {
  fofsi := finaloffsetintermediate{}

  fofsi.tilepos = make([]int, 0, 1024)

  var dummy uint64
  var dn int

  n:=0

  // nrecord
  dummy = byte2uint64(b[n:n+8])
  n += 8

  nrec := int(dummy)

  // data record byte length
  dummy = byte2uint64(b[n:n+8])
  n += 8

  bytelen := int(dummy)

  code := make([]byte, nrec)

  for i:=0 ;i<nrec; i++ {
    code[i] = b[n+i]
  }
  n+=nrec

  pos:=nrec
  for pos<bytelen {
    dummy,dn = dlug.ConvertUint64(b[n:])
    n+=dn
    pos += dn


    tilestep := int(dummy)
    fofsi.variant_ints = append(fofsi.variant_ints, tilestep)
    fofsi.tilepos = append(fofsi.tilepos, tilestep)

    dummy,dn = dlug.ConvertUint64(b[n:])
    n+=dn
    pos += dn

    nallele := int(dummy)
    fofsi.variant_ints = append(fofsi.variant_ints, nallele)

    for i:=0; i<nallele; i++ {
      dummy,dn = dlug.ConvertUint64(b[n:])
      n+=dn
      pos += dn

      len_allele_knot := int(dummy)
      fofsi.variant_ints = append(fofsi.variant_ints, len_allele_knot)

      for j:=0; j<len_allele_knot; j++ {

        dummy,dn = dlug.ConvertUint64(b[n:])
        n+=dn
        pos += dn

        var_id := int(dummy)
        fofsi.variant_ints = append(fofsi.variant_ints, var_id)

        dummy,dn = dlug.ConvertUint64(b[n:])
        n+=dn
        pos += dn

        span := int(dummy)
        fofsi.variant_ints = append(fofsi.variant_ints, span)

      }

    }

  }

  return fofsi

}

func bytes_from_finaloffsetintermediate(fofsi finaloffsetintermediate) []byte {
  buf := make([]byte, 64)
  fof_bytes := make([]byte, 0, 1024)

  // Number of records
  tobyte64(buf, uint64(len(fofsi.tilepos)))
  fof_bytes = append(fof_bytes, buf[0:8]...)

  // redundant...
  code := make([]byte, len(fofsi.tilepos))
  for i:=0; i<len(code); i++ { code[i] = 0 }

  data_bytes := make([]byte, 0, 1024)
  for i:=0; i<len(fofsi.variant_ints); i++ {
    vbytes := dlug.MarshalUint64(uint64(fofsi.variant_ints[i]))
    data_bytes = append(data_bytes, vbytes...)
  }

  // byte length of data record 
  //
  bytecount := uint64(len(code) + len(data_bytes))
  tobyte64(buf, bytecount)
  fof_bytes = append(fof_bytes, buf[0:8]...)

  // code section
  //
  fof_bytes = append(fof_bytes, code...)

  // data records
  //
  fof_bytes = append(fof_bytes, data_bytes...)

  return fof_bytes
}

func construct_final_offset_intermediate(ctx *CGFContext, prep_vector []vectorelement) finaloffsetintermediate {
  fofsi := finaloffsetintermediate{}

  fofsi.tilepos = make([]int, 0, 1024)
  fofsi.variant_ints = make([]int, 0, 1024)

  for i:=0; i<len(prep_vector); i++ {
    if prep_vector[i].fin_ovf_flag {
      fofsi.tilepos = append(fofsi.tilepos, i)

      knot := prep_vector[i].knot
      fofsi.variant_ints = append(fofsi.variant_ints,i)
      fofsi.variant_ints = append(fofsi.variant_ints, 2)
      for allele:=0 ; allele<2; allele++ {
        fofsi.variant_ints = append(fofsi.variant_ints, len(knot.varid[allele]))
        for i:=0; i<len(knot.varid[allele]); i++ {
          fofsi.variant_ints = append(fofsi.variant_ints, knot.varid[allele][i])
          fofsi.variant_ints = append(fofsi.variant_ints, knot.span[allele][i])
        }
      }

    }
  }

  return fofsi
}

func construct_uint64_vector(ctx *CGFContext, prep_vector []vectorelement) []uint64 {

  ret_vec := make([]uint64, 0, 1024)

  for i:=0; i<len(prep_vector); i+=32 {
    var cur_v uint64

    m := 32

    if i+32 > len(prep_vector) { m = len(prep_vector)%32 }

    for j:=0; j<m; j++ {

      if prep_vector[i+j].canon_flag { continue }

      cur_v |= (1<<(32+uint(j)))

      if prep_vector[i+j].cache_flag {

        // no hexit to set
        //
        if prep_vector[i+j].ovf_cache_flag { continue }

        // span is 0 hexit
        //
        if prep_vector[i+j].span_flag { continue }


        // generic overflow
        //
        if prep_vector[i+j].ovf_flag {
          cur_v |= 0xf << (4*uint(prep_vector[i+j].hexit_pos))
          continue
        }

        cur_v |= uint64( (uint(prep_vector[i+j].knot.tilemap_pos) & 0xf) << (4*uint(prep_vector[i+j].hexit_pos)) )

      }

    }

    ret_vec = append(ret_vec, cur_v)
    cur_v=0
  }

  return ret_vec

}

//====   _             _       _                               _ _       _
//====  | | ___   __ _(_)_ __ | |_ ___ _ __ _ __ ___   ___  __| (_) __ _| |_ ___
//====  | |/ _ \ / _` | | '_ \| __/ _ \ '__| '_ ` _ \ / _ \/ _` | |/ _` | __/ _ \
//====  | | (_) | (_| | | | | | ||  __/ |  | | | | | |  __/ (_| | | (_| | ||  __/
//====  |_|\___/ \__, |_|_| |_|\__\___|_|  |_| |_| |_|\___|\__,_|_|\__,_|\__\___|
//====              |_|

func loqintermediate_cmp(loqi0, loqi1 loqintermediate) error {
  if loqi0.count != loqi1.count {
    return fmt.Errorf( fmt.Sprintf("count mismatch: %d != %d", loqi0.count, loqi1.count) )
  }

  if loqi0.code != loqi1.code {
    return fmt.Errorf( fmt.Sprintf("code mismatch: %d != %d", loqi0.code, loqi1.code) )
  }

  if loqi0.stride != loqi1.stride {
    return fmt.Errorf( fmt.Sprintf("stride mismatch: %d != %d", loqi0.stride, loqi1.stride) )
  }

  if len(loqi0.tilepos) != len(loqi1.tilepos) {
    return fmt.Errorf( fmt.Sprintf("tilepos length mismatch: %d != %d", len(loqi0.tilepos), len(loqi1.tilepos)) )
  }

  if len(loqi0.homflag) != len(loqi1.homflag) {
    return fmt.Errorf( fmt.Sprintf("homflag length mismatch: %d != %d", len(loqi0.homflag), len(loqi1.homflag)) )
  }

  for i:=0; i<len(loqi0.homflag); i++ {
    if loqi0.homflag[i] != loqi1.homflag[i] {
      return fmt.Errorf( fmt.Sprintf("homflag mismatch at %d: %d != %d", i, loqi0.homflag[i], loqi1.homflag[i]) )
    }
  }

  if len(loqi0.loqinfo_ints) != len(loqi1.loqinfo_ints) {
    return fmt.Errorf( fmt.Sprintf("loqinfo_ints length mismatch: %d != %d", len(loqi0.loqinfo_ints), len(loqi1.loqinfo_ints)) )
  }

  for i:=0; i<len(loqi0.loqinfo_ints); i++ {
    if loqi0.loqinfo_ints[i] != loqi1.loqinfo_ints[i] {
      return fmt.Errorf( fmt.Sprintf("loqinfo_ints mismatch at %d: %d != %d", i, loqi0.loqinfo_ints[i], loqi1.loqinfo_ints[i]) )
    }
  }

  // we don't care about trailing overflow from the byte
  //
  if ((len(loqi0.loq_flag)+7)/8) != ((len(loqi1.loq_flag)+7)/8) {
    return fmt.Errorf( fmt.Sprintf("loq_flag length mismatch: %d != %d", len(loqi0.loq_flag), len(loqi1.loq_flag)) )
  }

  mm := len(loqi0.loq_flag)
  if len(loqi1.loq_flag) < mm { mm = len(loqi1.loq_flag) }
  for i:=0; i<mm; i++ {
    if loqi0.loq_flag[i] != loqi1.loq_flag[i] {
      return fmt.Errorf( fmt.Sprintf("loq_flag mismatch at %d: %v != %v", i, loqi0.loq_flag[i], loqi1.loq_flag[i]) )
    }
  }

  return nil

}

func loqintermediate_from_bytes(b []byte) loqintermediate {
  loqi := loqintermediate{}

  var dummy uint64
  var dn int ; _ = dn

  n:=0

  dummy = byte2uint64(b[n:n+8])
  n+=8

  rec_count := int(dummy)
  loqi.count = rec_count

  //DEBUG
  fmt.Printf("rec_count %d\n", rec_count)

  dummy = byte2uint64(b[n:n+8])
  n+=8

  code := int(dummy)
  loqi.code = code

  //DEBUG
  fmt.Printf("code: %d\n", code)

  dummy = byte2uint64(b[n:n+8])
  n+=8

  stride := int(dummy)
  loqi.stride = stride

  //DEBUG
  fmt.Printf("stride: %d\n", stride)

  offset_idx := make([]int, (rec_count+stride-1)/stride)
  for i:=0; i<(rec_count+stride-1)/stride; i++ {
    dummy = byte2uint64(b[n:n+8])
    n+=8

    offset_idx[i] = int(dummy)
  }

  //loqi.offset_idx = append(loqi.offset_idx, offset_idx...)

  //DEBUG
  for i:=0; i<len(offset_idx); i++ {
    fmt.Printf("offset[%d]: %d\n", i, offset_idx[i])
  }




  tilepos_idx := make([]int, (rec_count+stride-1)/stride)
  for i:=0; i<(rec_count+stride-1)/stride; i++ {
    dummy = byte2uint64(b[n:n+8])
    n+=8

    tilepos_idx[i] = int(dummy)
  }


  //loqi.tilepos_idx

  //DEBUG
  for i:=0; i<len(tilepos_idx); i++ {
    fmt.Printf("tilepos[%d]: %d\n", i, tilepos_idx[i])
  }


  homflag := make([]byte, (rec_count+7)/8)
  for i:=0; i<(rec_count+7)/8; i++ {
    homflag[i] = b[n]
    n++
  }

  for i:=0; i<rec_count; i++ {
    tf := false
    if (homflag[i/8] & (1<<uint(i%8)))!=0 { tf = true }
    loqi.homflag = append(loqi.homflag, tf)
  }

  cur_tilepos := 0
  for i:=0; i<rec_count; i++ {
    if (i%loqi.stride) == 0 {
      cur_tilepos = tilepos_idx[i/loqi.stride]
    }
    loqi.tilepos = append(loqi.tilepos, cur_tilepos)
  }



  // loq flag size on vector
  //
  dummy = byte2uint64(b[n:n+8])
  n+=8

  loqflag_bytecount := int(dummy)

  fmt.Printf("loqflag_bytecount %d\n", loqflag_bytecount)

  loq_flag_vec := b[n:n+loqflag_bytecount]
  n+=loqflag_bytecount

  for i:=0; i<8*loqflag_bytecount; i++ {
    tf := false
    if (loq_flag_vec[i/8] & (1<<uint(i%8))) != 0 { tf = true }
    loqi.loq_flag = append(loqi.loq_flag, tf)
  }


  // size of loq array
  //
  dummy = byte2uint64(b[n:n+8])
  n+=8

  loq_info_byte_count := int(dummy) ; _ = loq_info_byte_count



  //DEBUG
  fmt.Printf("loq byte count %d\n", loq_info_byte_count)
  fmt.Printf("BYTE LOQ DEBUG\n")

  // man loq array
  //
  rec_pos:=0
  byte_offset := 0
  for byte_offset < loq_info_byte_count {

    ntile := make([]int, 1)

    dummy,dn := dlug.ConvertUint64(b[n:])
    n+=dn
    byte_offset+=dn

    ntile[0] = int(dummy)
    loqi.loqinfo_ints = append(loqi.loqinfo_ints, int(ntile[0]))

    //fmt.Printf("rec_pos %d, loqi.homflag len %d\n", rec_pos, len(loqi.homflag))
    fmt.Printf(" loq[%d] n[0] %d\n", rec_pos, ntile[0])

    if !loqi.homflag[rec_pos] {


      dummy,dn := dlug.ConvertUint64(b[n:])
      n+=dn
      byte_offset+=dn

      ntile = append(ntile, int(dummy))

      fmt.Printf("+loq[%d] n[1] %d\n", rec_pos, ntile[1])

      loqi.loqinfo_ints = append(loqi.loqinfo_ints, int(ntile[1]))
    }

    fmt.Printf(" %v\n", ntile)

    for allele:=0; allele<len(ntile); allele++ {

      for i:=0; i<ntile[allele]; i++ {


        dummy,dn := dlug.ConvertUint64(b[n:])
        n+=dn
        byte_offset+=dn

        m := int(dummy)
        loqi.loqinfo_ints = append(loqi.loqinfo_ints, int(m))

        fmt.Printf("   loq[%d][%d] mlen %d\n", rec_pos, allele, m)

        //run_sum := 0
        for j:=0; j<m; j+=2 {
          dummy,dn := dlug.ConvertUint64(b[n:])
          n+=dn
          byte_offset+=dn

          delpos:=int(dummy)

          dummy,dn = dlug.ConvertUint64(b[n:])
          n+=dn
          byte_offset+=dn

          l:=int(dummy)

          fmt.Printf("      loq[%d][%d] %d+%d\n", rec_pos, allele, delpos, l)

          //loqi.loqinfo_ints = append(loqi.loqinfo_ints, run_sum+delpos)
          loqi.loqinfo_ints = append(loqi.loqinfo_ints, delpos)
          loqi.loqinfo_ints = append(loqi.loqinfo_ints, l)
          //run_sum += delpos
        }
      }
    }

    rec_pos++

  }

  return loqi

}

func bytes_from_loqintermediate(loqi loqintermediate) []byte {
  buf := make([]byte, 64)
  loq_bytes := make([]byte, 0, 1024)

  tobyte64(buf, uint64(loqi.count))
  loq_bytes = append(loq_bytes, buf[0:8]...)

  tobyte64(buf, uint64(loqi.code))
  loq_bytes = append(loq_bytes, buf[0:8]...)

  tobyte64(buf, uint64(loqi.stride))
  loq_bytes = append(loq_bytes, buf[0:8]...)

  offset_idx := make([]uint64, 0, 1024)
  tilepos_idx := make([]uint64, 0, 1024)

  loqinfo_bytes := make([]byte, 0, 1024)

  loq_flag := make([]byte, (len(loqi.loq_flag) + 7)/8)
  for i:=0; i<len(loqi.loq_flag); i++ {
    if loqi.loq_flag[i] { loq_flag[i/8] |= 1<<uint(i%8) }
  }

  cur_rec := 0
  byte_offset := 0 ; _ = byte_offset
  p:=0
  for p<len(loqi.loqinfo_ints) {

    if (cur_rec%loqi.stride) == 0 {
      offset_idx = append(offset_idx, uint64(len(loqinfo_bytes)))
      tilepos_idx = append(tilepos_idx, uint64(loqi.tilepos[cur_rec]))
    }

    ntile := make([]int, 1)

    ma := dlug.MarshalUint64(uint64(loqi.loqinfo_ints[p]))
    loqinfo_bytes = append(loqinfo_bytes, ma...)

    ntile[0]=loqi.loqinfo_ints[p]
    p++

    if !loqi.homflag[cur_rec] {

      mb := dlug.MarshalUint64(uint64(loqi.loqinfo_ints[p]))
      loqinfo_bytes = append(loqinfo_bytes, mb...)

      ntile = append(ntile, loqi.loqinfo_ints[p])
      p++
    }

    for allele:=0; allele<len(ntile); allele++ {
      for i:=0; i<ntile[allele]; i++ {

        m:=loqi.loqinfo_ints[p]
        p++

        mb := dlug.MarshalUint64(uint64(m))
        loqinfo_bytes = append(loqinfo_bytes, mb...)

        for j:=0; j<m; j+=2 {
          pos := loqi.loqinfo_ints[p]
          p++

          l := loqi.loqinfo_ints[p]
          p++

          x := dlug.MarshalUint64(uint64(pos))
          loqinfo_bytes = append(loqinfo_bytes, x...)

          y := dlug.MarshalUint64(uint64(l))
          loqinfo_bytes = append(loqinfo_bytes, y...)
        }

      }

    }

    cur_rec++

  }

  loq_info_byte_count := len(loqinfo_bytes)

  homflag := make([]byte, (cur_rec+7)/8)
  for i:=0; i<cur_rec; i++ {

    if loqi.homflag[i] {
      homflag[i/8] |= 1<<uint(i%8)
    }
  }

  for i:=0; i<len(offset_idx); i++ {
    tobyte64(buf, offset_idx[i])
    loq_bytes = append(loq_bytes, buf[0:8]...)
  }

  for i:=0; i<len(tilepos_idx); i++ {
    tobyte64(buf, tilepos_idx[i])
    loq_bytes = append(loq_bytes, buf[0:8]...)
  }

  loq_bytes = append(loq_bytes, homflag...)

  tobyte64(buf, uint64(len(loq_flag)))
  loq_bytes = append(loq_bytes, buf[0:8]...)
  loq_bytes = append(loq_bytes, loq_flag...)

  fmt.Printf("writing loq_flag %d %v\n", len(loq_flag), buf[0:8])

  tobyte64(buf, uint64(loq_info_byte_count))
  loq_bytes = append(loq_bytes, buf[0:8]...)
  loq_bytes = append(loq_bytes, loqinfo_bytes...)

  return loq_bytes
}

func _loq_hom(nocall_start_len [][][]int) bool {
  if len(nocall_start_len)==1 { return true }
  if len(nocall_start_len)!=2 { return false }

  a := nocall_start_len[0]
  b := nocall_start_len[1]

  if len(a)!=len(b) { return false }
  for i:=0; i<len(a); i++ {
    if len(a[i]) != len(b[i]) { return false }
    for j:=0; j<len(a[i]); j++ {
      if a[i][j]!=b[i][j] { return false }
    }
  }

  return true
}

func construct_loq_intermediate(ctx *CGFContext, prep_vector []vectorelement) loqintermediate {
  loqi := loqintermediate{}

  loqi.code = 0
  loqi.stride = 256

  // fill out loq bit vector
  //
  for i:=0; i<len(prep_vector); i++ {

    if prep_vector[i].span_flag {
      loqi.loq_flag = append(loqi.loq_flag, false)
      continue
    }

    if prep_vector[i].loq_flag {
      loqi.loq_flag = append(loqi.loq_flag, true)
    } else {
      loqi.loq_flag = append(loqi.loq_flag, false)
    }

  }

  loqi.count = 0

  // populate loqinfo_ints
  //
  for i:=0; i<len(prep_vector); i++ {
    if prep_vector[i].loq_flag {
      loqi.tilepos = append(loqi.tilepos, i)
      loqi.count++

      nocall_start_len := prep_vector[i].knot.nocall_start_len

      if _loq_hom(nocall_start_len) {

        // Hom
        //

        loqi.homflag = append(loqi.homflag, true)
        loqi.loqinfo_ints = append(loqi.loqinfo_ints, len(nocall_start_len[0]))


        //DEBUG
        fmt.Printf("%x - ** N %d\n", i, len(nocall_start_len[0]))



        for ii:=0; ii<len(nocall_start_len[0]); ii++ {
          loqi.loqinfo_ints = append(loqi.loqinfo_ints, len(nocall_start_len[0][ii]))

          //DEBUG
          fmt.Printf("%x -  ** m(%d) %d\n", i, ii, len(nocall_start_len[0][ii]))


          start := 0
          for jj:=0; jj<len(nocall_start_len[0][ii]); jj+=2 {
            loqi.loqinfo_ints = append(loqi.loqinfo_ints, nocall_start_len[0][ii][jj]-start)
            loqi.loqinfo_ints = append(loqi.loqinfo_ints, nocall_start_len[0][ii][jj+1])
            start = nocall_start_len[0][ii][jj]

            //DEBUG
            fmt.Printf("%x -  ** %d+%d\n", i, nocall_start_len[0][ii][jj], nocall_start_len[0][ii][jj+1])


          }

        }


      } else {

        // Het
        //
        loqi.homflag = append(loqi.homflag, false)
        loqi.loqinfo_ints = append(loqi.loqinfo_ints, len(nocall_start_len[0]))
        loqi.loqinfo_ints = append(loqi.loqinfo_ints, len(nocall_start_len[1]))

        //DEBUG
        fmt.Printf("%x - ** N %d (bonk)\n", i, len(nocall_start_len[0]))
        fmt.Printf("%x - ** N %d (bonk)\n", i, len(nocall_start_len[1]))

        for allele:=0; allele<2; allele++ {
          for ii:=0; ii<len(nocall_start_len[allele]); ii++ {
            loqi.loqinfo_ints = append(loqi.loqinfo_ints, len(nocall_start_len[allele][ii]))

            //DEBUG
            fmt.Printf("%x -  ** m(%d) %d\n", i, ii, len(nocall_start_len[allele][ii]))


            start := 0
            for jj:=0; jj<len(nocall_start_len[allele][ii]); jj+=2 {
              loqi.loqinfo_ints = append(loqi.loqinfo_ints, nocall_start_len[allele][ii][jj]-start)
              loqi.loqinfo_ints = append(loqi.loqinfo_ints, nocall_start_len[allele][ii][jj+1])
              start = nocall_start_len[allele][ii][jj]

              //DEBUG
              fmt.Printf("%x -  ** %d+%d\n", i, nocall_start_len[allele][ii][jj], nocall_start_len[allele][ii][jj+1])

            }

          }
        }

      }

    }
  }


  return loqi
}

//====                 _ _
//====    ___ _ __ ___ (_) |_
//====   / _ \ '_ ` _ \| | __|
//====  |  __/ | | | | | | |_
//====   \___|_| |_| |_|_|\__|


func emit_intermediate(ctx *CGFContext, path_idx int, allele_path [][]TileInfo) error {
  debug_output:=true

  cgf := ctx.CGF ; _ = cgf
  sglf := ctx.SGLF

  span_sum := 0
  step_idx0,step_idx1 := 0,0

  knot := cgfintermediate{}
  _init_knot(&knot)

  tileKnot := make([]cgfintermediate, 0, 1024)

  // Construct the intermediate string of knots
  //
  for (step_idx0<len(allele_path[0])) || (step_idx1<len(allele_path[1])) {

    if span_sum >= 0 {
      s,e := _add_knot(&knot, 0, step_idx0, allele_path[0][step_idx0], sglf)
      if e!=nil { panic(e) }

      step_idx0++
      span_sum -= s
    } else {
      s,e := _add_knot(&knot, 1, step_idx1, allele_path[1][step_idx1], sglf)
      if e!=nil { panic(e) }

      step_idx1++
      span_sum += s
    }

    if span_sum==0 {

      _knot_tot_span(&knot)
      knot.tilemap_key = create_tilemap_string_lookup2(knot.varid[0], knot.span[0], knot.varid[1], knot.span[1])
      tileKnot = append(tileKnot, knot)

      knot = cgfintermediate{}
      _init_knot(&knot)
    }

  }

  // Prep for binary representation
  //
  prep_vector := make([]vectorelement, 0, 1024)
  cache_counter := 0

  for ind:=0; ind<len(tileKnot); ind++ {
    knot := tileKnot[ind]

    cur_step := knot.step[0][0]
    next_step := cur_step + knot.tot_span

    if (cur_step%(32))==0 {
      cache_counter = 0
    }

    vec_ele := vectorelement{}

    // We have a canonical tile.  Add it and move on
    //
    if (!knot.loq_flag) && (knot.tot_span == 1) && (knot.varid[0][0] == 0) && (knot.varid[1][0]==0) {
      vec_ele.canon_flag = true
      vec_ele.knot = knot
      prep_vector = append(prep_vector, vec_ele)
      continue
    }

    if cache_counter > (32/4) {
      vec_ele.ovf_cache_flag = true
    }
    cache_counter++

    if knot.loq_flag { vec_ele.loq_flag = true }

    if _,ok := ctx.TileMapLookup[knot.tilemap_key] ; ok {
      // We've found it in the TileMap.  We can either
      // put it into the vector cache or we can put it into
      // the overflow table.  If it's either low quality
      // or more than 0xd, it goes into the overflow table.
      // Otherwise it can go in the vector cache.
      //

      //knot.tilemap_pos = tilemap_pos.TileMap
      knot.tilemap_pos = ctx.TileMapPosition[knot.tilemap_key]

      if knot.tilemap_pos >= 0xd { vec_ele.ovf_flag = true }
      if vec_ele.loq_flag { vec_ele.ovf_flag = true }

      if cache_counter > (32/4) {
        vec_ele.ovf_cache_flag = true
        vec_ele.ovf_flag = true
      } else {
        vec_ele.cache_flag = true
        vec_ele.hexit_pos = cache_counter-1
      }



      // If our cache can still hold hexits and our tilemap
      // entry isn't too big and it's high quality.
      //
      //if !vec_ele.ovf_cache_flag && !vec_ele.ovf_flag {
      //  vec_ele.cache_flag = true
      //  vec_ele.hexit_pos = cache_counter
      //}

    } else {

      // We couldn't find it in the TileMap table, so
      // we store it in the FinalOverflowMap table
      //

      vec_ele.fin_ovf_flag = true
      vec_ele.ovf_flag = true
      //if !vec_ele.ovf_cache_flag && !vec_ele.ovf_flag {
      if !vec_ele.ovf_cache_flag {

        if cache_counter > (32/4) {
          vec_ele.ovf_cache_flag = true
        } else {
          vec_ele.cache_flag = true
          vec_ele.hexit_pos = cache_counter-1
        }
      }

    }

    vec_ele.knot = knot
    prep_vector = append(prep_vector, vec_ele)

    // Add an explicit entry for spanning tiles
    //
    cur_step++
    for ; cur_step<next_step; cur_step++ {
      if (cur_step%(32))==0 {
        cache_counter = 0
      }
      cache_counter++

      span_vec_ele := vectorelement{}
      span_vec_ele.canon_flag = false
      span_vec_ele.span_flag = true

      if cache_counter > (32/4) {
        span_vec_ele.ovf_cache_flag = true
        span_vec_ele.ovf_flag = true
      } else {
        span_vec_ele.cache_flag = true
        span_vec_ele.hexit_pos = cache_counter-1
      }

      prep_vector = append(prep_vector, span_vec_ele)
    }

  }


  // ======================================================
  // ======================================================
  // ======================================================
  // ======================================================
  // ======================================================
  // ======================================================
  // ======================================================
  // ======================================================
  // ======================================================
  // ======================================================

  if debug_output {
    for i:=0; i<len(prep_vector); i++ {

      if (i%32) == 0 {
        fmt.Printf("#    con(^) ca(c) oca(>) ovf(/) fin(!) sp(~) lq(*) hexit vec ...\n")
      }

      if prep_vector[i].span_flag {
        //fmt.Printf("[%d+_] cache %v, ovf_cache %v\n", i, prep_vector[i].cache_flag, prep_vector[i].ovf_cache_flag)
        fmt.Printf("[%4x+_] ", i)
        //fmt.Printf(" con %v ca %v oca %v ovf %v fin %v sp %v lq %v x %v v %v ",
        fmt.Printf(" %v %v %v %v %v %v %v %v %v tmap(%d) ",
          _tf_(prep_vector[i].canon_flag, "^"),
          _tf_(prep_vector[i].cache_flag, "c"),
          _tf_(prep_vector[i].ovf_cache_flag, ">"),
          _tf_(prep_vector[i].ovf_flag, "/"),
          _tf_(prep_vector[i].fin_ovf_flag, "!"),
          _tf_(prep_vector[i].span_flag, "~"),
          _tf_(prep_vector[i].loq_flag, "*"),
          prep_vector[i].hexit_pos,
          prep_vector[i].vec_pos,
          prep_vector[i].knot.tilemap_pos)

        fmt.Printf("\n")

      } else {
        knot := prep_vector[i].knot
        fmt.Printf("[%4x+%x] ", knot.step[0][0], knot.tot_span)

        if prep_vector[i].cache_flag && (prep_vector[i].hexit_pos==0) {
          fmt.Printf(" %v %v %v %v %v %v %v %v %v tmap(%d) ",
            _tf_(prep_vector[i].canon_flag, "^"),
            _tf_(prep_vector[i].cache_flag, "c"),
            _tf_(prep_vector[i].ovf_cache_flag, ">"),
            _tf_(prep_vector[i].ovf_flag, "/"),
            _tf_(prep_vector[i].fin_ovf_flag, "!"),
            _tf_(prep_vector[i].span_flag, "~"),
            _tf_(prep_vector[i].loq_flag, "*"),
            "$",
            prep_vector[i].vec_pos,
            prep_vector[i].knot.tilemap_pos)

        } else {
          fmt.Printf(" %v %v %v %v %v %v %v %v %v tmap(%d) ",
            _tf_(prep_vector[i].canon_flag, "^"),
            _tf_(prep_vector[i].cache_flag, "c"),
            _tf_(prep_vector[i].ovf_cache_flag, ">"),
            _tf_(prep_vector[i].ovf_flag, "/"),
            _tf_(prep_vector[i].fin_ovf_flag, "!"),
            _tf_(prep_vector[i].span_flag, "~"),
            _tf_(prep_vector[i].loq_flag, "*"),
            prep_vector[i].hexit_pos,
            prep_vector[i].vec_pos,
            prep_vector[i].knot.tilemap_pos)
        }

        for allele:=0; allele<2; allele++ {
          if allele==0 {
            fmt.Printf("A(");
          } else {
            fmt.Printf(") : B(");
          }
          for j:=0; j<len(knot.varid[allele]); j++ {
            if j>0 { fmt.Printf(":") }
            ch := "_"
            if knot.loq[allele][j] { ch = "*" }
            fmt.Printf("%x%s%x+%x",
              knot.step[allele][j],
              ch,
              knot.varid[allele][j],
              knot.span[allele][j])
          }
        }
        fmt.Printf(")")

        fmt.Printf("\n")

      }

    }
  }

  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================

  vec64 := construct_uint64_vector(ctx, prep_vector)

  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================

  if debug_output {
    vec_bytes := make([]byte, 8*len(vec64))
    for i:=0; i<len(vec64); i++ {
      tobyte64(vec_bytes[8*i:], vec64[i])
    }

    fmt.Printf(">>>returned vec %d\n", len(vec64))
    for i:=0; i<len(vec64); i++ {
      if (i%4)==0 { fmt.Printf("\n") }
      fmt.Printf("[%3x,%4x] %8x.%8x |", i, 32*i, (vec64[i]&(0xffffffff<<32))>>32, vec64[i]&0xffffffff)
    }
    fmt.Printf("\n")


    random_start := 0x12af
    random_n := 0x120
    random_ovf_count := 0
    for i:=random_start; i<(random_start+random_n); i++ {
      if prep_vector[i].ovf_flag || prep_vector[i].ovf_cache_flag {
        random_ovf_count++
      }
    }

    check_ovf_count := CountOverflowVectorUint64(vec64, random_start, random_start+random_n)

    fmt.Printf("CHECKING (step %x+%x(%x)) real:%d check:%d\n", random_start, random_n, random_start+random_n, random_ovf_count, check_ovf_count)


    if debug_output {

      for i:=0; i<len(tileKnot); i++ {
        fmt.Printf("[%4x+%x] ", tileKnot[i].step[0][0], tileKnot[i].tot_span)

        for allele:=0; allele<2; allele++ {
          if allele==0 {
            fmt.Printf("A(");
          } else {
            fmt.Printf(") : B(");
          }
          for j:=0; j<len(tileKnot[i].varid[allele]); j++ {
            if j>0 { fmt.Printf(":") }
            ch := "_"
            if tileKnot[i].loq[allele][j] { ch = "*" }
            fmt.Printf("%x%s%x+%x",
              tileKnot[i].step[allele][j],
              ch,
              tileKnot[i].varid[allele][j],
              tileKnot[i].span[allele][j])
          }
        }
        fmt.Printf(")\n")

      }


      fmt.Printf("LOQ INFO (tileKnot)\n")
      for i:=0; i<len(tileKnot); i++ {
        fmt.Printf("[%4x+%x] loq ", tileKnot[i].step[0][0], tileKnot[i].tot_span)

        for allele:=0; allele<2; allele++ {
          if allele==0 {
            fmt.Printf("A(");
          } else {
            fmt.Printf(") : B(");
          }
          for j:=0; j<len(tileKnot[i].nocall_start_len[allele]); j++ {
            if j>0 { fmt.Printf(",") }

            fmt.Printf("{%d}", j)
            for k:=0; k<len(tileKnot[i].nocall_start_len[allele][j]); k+=2 {
              fmt.Printf(";%d+%d",
                tileKnot[i].nocall_start_len[allele][j][k],
                tileKnot[i].nocall_start_len[allele][j][k+1])
            }
          }
        }
        fmt.Printf(")\n")

      }

    }


  }


  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================


  ofsi := construct_offset_intermediate(ctx, prep_vector)


  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================


  if debug_output {

    fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n")

    fmt.Printf("BYTES FROM ofsi to byte_ofsi_test0\n")
    byte_ofsi_test0 := bytes_from_offsetintermediate(ofsi) ; _ = byte_ofsi_test0


    fmt.Printf("BYTES TO byte_ofsi_test0 to ofsi_throwaway\n")
    ofsi_throwaway := offsetintermediate_from_bytes(byte_ofsi_test0) ; _ = ofsi_throwaway

    for i:=0; i<len(ofsi_throwaway.tilepos_idx); i++ {
      fmt.Printf("  hrm: ofsi_throwaway.tilepos_idx[%d] %d\n", i, ofsi_throwaway.tilepos_idx[i])
    }

    err := offsetintermediate_cmp(ofsi_throwaway, ofsi)
    fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> %v\n", err)

    fmt.Printf("BYTES FROM ofsi_throwaway to byte_ofsi_test1\n")
    byte_ofsi_test1 := bytes_from_offsetintermediate(ofsi_throwaway)
    _ = byte_ofsi_test1


    fmt.Printf("BYTES TO byte_ofsi_test1 to ofsi_throwaway1\n")
    ofsi_throwaway1 := offsetintermediate_from_bytes(byte_ofsi_test0) ; _ = ofsi_throwaway1

    /*
    for i:=0; i<len(ofsi_throwaway.offset_idx); i++ {
      fmt.Printf("  offset_idx[%d] %d %d\n", i, ofsi_throwaway.offset_idx[i], ofsi_throwaway1.offset_idx[i])
    }

    for i:=0; i<len(ofsi_throwaway.offset_idx); i++ {
      fmt.Printf("  tilepos_idx[%d] %d %d\n", i, ofsi_throwaway.tilepos_idx[i], ofsi_throwaway1.tilepos_idx[i])
    }
    */


    fmt.Printf("!!! %d %d\n", len(byte_ofsi_test0), len(byte_ofsi_test1))
    for i:=0; i<len(byte_ofsi_test0); i++ {
      if byte_ofsi_test0[i] != byte_ofsi_test1[i] {
        fmt.Printf("mismatched byte at %d (%d != %d)\n", i, byte_ofsi_test0[i], byte_ofsi_test1[i])
      }
    }

    fmt.Printf("OFSI (%d (%x))\n", len(ofsi.tilemap), len(ofsi.tilemap))
    for i:=0; i<len(ofsi.tilemap); i++ {
      fmt.Printf("[%d] {%x} %d (%x) fovf?%v span?%v\n", i, ofsi.tilepos[i], ofsi.tilemap[i], ofsi.tilemap[i], ofsi.final_overflow_flag[i], ofsi.span_flag[i])
    }

    fmt.Printf("OFSI CHECK\n")
    stride_tile_pos := 0
    stride_ofs := 0
    for i:=0; i<len(ofsi.tilemap); i++ {

      if i%256 == 0 {
        stride_tile_pos = ofsi.tilepos[i]
        stride_ofs = i
      }

      check_ovf_count := CountOverflowVectorUint64(vec64, stride_tile_pos, ofsi.tilepos[i])
      check_ovf_count += stride_ofs

      fmt.Printf("[%d] {%x} %d (%x) ==> ovf_count: %d (%x)\n", i, ofsi.tilepos[i], ofsi.tilemap[i], ofsi.tilemap[i], check_ovf_count, check_ovf_count)
      if check_ovf_count != i {
        real_ovf_count := vectorelement_ovf_count(prep_vector, stride_tile_pos, ofsi.tilepos[i] - stride_tile_pos) 
        fmt.Printf("ERROR!!!! check_ovf_count %d != i %d (real %d)\n", check_ovf_count, i, real_ovf_count)


      }


    }

  }


  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================

  fofsi := construct_final_offset_intermediate(ctx, prep_vector)

  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================


  if debug_output {


    fmt.Printf("FOFSI (%d (%x))\n", len(fofsi.tilepos), len(fofsi.tilepos))

    cur_int_pos := 0

    for i:=0; i<len(fofsi.tilepos); i++ {
      fmt.Printf("[%d] tilepos{%x}", i, fofsi.tilepos[i])

      step := fofsi.variant_ints[cur_int_pos] ; cur_int_pos++
      nallele := fofsi.variant_ints[cur_int_pos] ; cur_int_pos++

      fmt.Printf(" step:%x N(%d)", step, nallele)

      n_a_allele := fofsi.variant_ints[cur_int_pos] ; cur_int_pos++
      fmt.Printf(" A(%d)[", n_a_allele)

      for ii:=0; ii<n_a_allele; ii++ {
        fmt.Printf(" %x+%x", fofsi.variant_ints[cur_int_pos], fofsi.variant_ints[cur_int_pos+1])
        cur_int_pos+=2
      }
      fmt.Printf(" ]")

      n_b_allele := fofsi.variant_ints[cur_int_pos] ; cur_int_pos++
      fmt.Printf(" B(%d)[", n_b_allele)

      for ii:=0; ii<n_b_allele; ii++ {
        fmt.Printf(" %x+%x", fofsi.variant_ints[cur_int_pos], fofsi.variant_ints[cur_int_pos+1])
        cur_int_pos+=2
      }
      fmt.Printf(" ]\n")


    }


    fofsi_bytes0 := bytes_from_finaloffsetintermediate(fofsi) ; _ = fofsi_bytes0
    fofsi_temp1 := finaloffsetintermediate_from_bytes(fofsi_bytes0) ; _ = fofsi_temp1
    fofsi_bytes1 := bytes_from_finaloffsetintermediate(fofsi_temp1)

    fofsi_cnv := finaloffsetintermediate_from_bytes(fofsi_bytes1)
    fmt.Printf("FOFSI BYTES %d %d\n", len(fofsi_bytes0), len(fofsi_bytes1))

    if len(fofsi_bytes0) != len(fofsi_bytes1) {
      fmt.Printf("ERROR: length mismatch for fofsi_bytes %d != %d\n", len(fofsi_bytes0), len(fofsi_bytes1))
    } else {
      for i:=0; i<len(fofsi_bytes0); i++ {
        if fofsi_bytes0[i] != fofsi_bytes1[i] {
          fmt.Printf("ERROR: byte mismatch for fofsi_bytes %d: %d != %d\n", i, fofsi_bytes0[i], fofsi_bytes1[i])
        }
      }
    }

    err := finaloffsetintermediate_cmp(fofsi, fofsi_cnv)
    if err!=nil { fmt.Printf("ERROR: %v\n", err) }

  }

  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================

  loqi := construct_loq_intermediate(ctx, prep_vector)
  _ = loqi

  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================
  //=====================================================

  if debug_output {

    fmt.Printf("LOQ DEBUG\n")
    fmt.Printf("loqi count %d\n", loqi.count)
    fmt.Printf("loqi stride %d\n", loqi.stride)

    p:=0
    for i:=0; i<len(loqi.tilepos); i++ {
      fmt.Printf("{%x} %v ", loqi.tilepos[i], loqi.homflag[i])

      if loqi.homflag[i] {

        n := loqi.loqinfo_ints[p]
        p++

        fmt.Printf("[%d]", n)

        for ii:=0; ii<n; ii++ {
          m := loqi.loqinfo_ints[p]
          p++

          fmt.Printf(" (%d)", m)

          st:=0
          for jj:=0; jj<m; jj+=2 {
            fmt.Printf(";%d+%d", loqi.loqinfo_ints[p]+st, loqi.loqinfo_ints[p+1])
            st += loqi.loqinfo_ints[p]
            p+=2
          }
        }

      } else {

        n0 := loqi.loqinfo_ints[p]
        p++
        n1 := loqi.loqinfo_ints[p]
        p++

        fmt.Printf(" (%d,%d)", n0,n1)

        for ii:=0; ii<n0; ii++ {
          m := loqi.loqinfo_ints[p]
          p++

          fmt.Printf(" (%d)", m)

          st:=0
          for jj:=0; jj<m; jj+=2 {
            fmt.Printf(";%d+%d", loqi.loqinfo_ints[p]+st, loqi.loqinfo_ints[p+1])
            st += loqi.loqinfo_ints[p]
            p+=2
          }
        }

        fmt.Printf(" :: ")

        for ii:=0; ii<n1; ii++ {
          m := loqi.loqinfo_ints[p]
          p++

          fmt.Printf(" (%d)", m)

          st:=0
          for jj:=0; jj<m; jj+=2 {
            fmt.Printf(";%d+%d", loqi.loqinfo_ints[p]+st, loqi.loqinfo_ints[p+1])
            st += loqi.loqinfo_ints[p]
            p+=2
          }
        }


      }

      fmt.Printf("\n")

    }


    loq_bytes0 := bytes_from_loqintermediate(loqi) ; _ = loq_bytes0
    loqi_test0 := loqintermediate_from_bytes(loq_bytes0)
    loq_bytes1 := bytes_from_loqintermediate(loqi_test0)

    fmt.Printf(">>> lOQ CMP: %v\n", loqintermediate_cmp(loqi, loqi_test0))
    fmt.Printf(">>>>>>>>>>>>>> LOQ FROM/TO BYTES %d %d\n", len(loq_bytes0), len(loq_bytes1))

    if len(loq_bytes0) != len(loq_bytes1) {
      fmt.Printf("ERROR: len(loq_bytes0) %d != len(loq_bytes1) %d\n", len(loq_bytes0), len(loq_bytes1))
    } else {
      for i:=0; i<len(loq_bytes0); i++ {
        if loq_bytes0[i] != loq_bytes1[i] {
          fmt.Printf("byte mismatch at %d: %d != %d\n", i, loq_bytes0[i], loq_bytes1[i])
        }
      }
    }


  }


  return nil

}

func _tf(b bool) string {
  if b { return "t" }
  return "."
}
func _tf_(b bool, s string) string {
  if b { return s }
  return "."
}
