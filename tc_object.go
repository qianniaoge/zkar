package javaserialize

type TCObject struct {
	ClassPointer *TCClassPointer
	ClassDatas   []*TCClassData
	Handler uint32
}

func (oo *TCObject) ToBytes() []byte {
	var bs = []byte{JAVA_TC_OBJECT}
	bs = append(bs, oo.ClassPointer.ToBytes()...)
	for _, data := range oo.ClassDatas {
		bs = append(bs, data.ToBytes()...)
	}

	return bs
}

func (oo *TCObject) ToString() string {
	var b = NewPrinter()
	b.Printf("TC_OBJECT - %s\n", Hexify(JAVA_TC_OBJECT))
	b.IncreaseIndent()
	b.Printf(oo.ClassPointer.ToString())
	b.Printf("@Handler - %v", oo.Handler)
	b.Printf("[]ClassData \n")
	b.IncreaseIndent()
	for _, data := range oo.ClassDatas {
		b.Printf(data.ToString())
		b.Printf("\n")
	}

	return b.String()
}

func readTCObject(stream *ObjectStream) (*TCObject, error) {
	var obj = new(TCObject)
	var err error

	_, _ = stream.ReadN(1)
	obj.ClassPointer, err = readTCClassPointer(stream)
	if err != nil {
		return nil, err
	}

	stream.AddReference(obj)
	if obj.ClassPointer.Flag == JAVA_TC_NULL {
		return obj, nil
	}

	bag, err := obj.ClassPointer.FindClassBag(stream)
	if err != nil {
		return nil, err
	}

	for i := len(bag.Classes) - 1; i >= 0; i-- {
		classData, err := readTCClassData(stream, bag.Classes[i])
		if err != nil {
			return nil, err
		}

		obj.ClassDatas = append(obj.ClassDatas, classData)
	}

	return obj, nil
}
