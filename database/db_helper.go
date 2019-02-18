package database

// GetImageAccountFromDB returns an ImageAccount if exists in the database
func GetImageAccountFromDB(hash string) (*ImageAccount, error) {
	image := &ImageAccount{}

	i, err := GetDB().Model(image).Get([]byte(hash))
	if err != nil {
		return nil, err
	}
	image = i.(*ImageAccount)
	return image, nil
}

// GetImageFromDB returns an ImageLvlDB if exists in the database
func GetImageFromDB(imgHash string) (*ImageLvlDB, error) {
	image := &ImageLvlDB{}
	i, err := GetDB().Model(image).Get([]byte(imgHash))
	if err != nil {
		return nil, err
	}
	image = i.(*ImageLvlDB)
	return image, nil
}
