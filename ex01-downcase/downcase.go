package downcase

func Downcase(line string) (string, error){
	s := make([]byte, len(line))

	for i := range s{
		symbol := line[i]
		if symbol >= 65 && symbol <=90 {
			symbol += 32;
		}
		s[i] = symbol
	}
	return string(s), nil
}
