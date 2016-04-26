package easylog

func GetLogger(filename string, level int) *Logger {
	w := NewFileWriter()
	pattern := ".%Y%M%D%H"
	w.SetFileName(filename)
	w.SetPathPattern(pattern)

	logger := NewLogger()
	if w == nil || logger == nil {
		return nil
	}
	logger.Register(w)
	logger.SetLevel(level)
	return logger
}
