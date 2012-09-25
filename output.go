package main


/*func BinIndex(index *Index) []byte {
	
}*/

func RepIndex(index *Index) string {
	s := "Index:\n"
	
	for i := 0; i < len(index.Sites); i++ {
		s += index.Sites[i].URL + "\n"
		for j := 0; j < len(index.Sites[i].Pages); j++ {
			s += "\t" + index.Sites[i].Pages[j].Path + "\n"
		}
	}
	
	return s + "\n\n"
}
