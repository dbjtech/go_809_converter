package helpers

type CarHelper struct {
	cache CacheHelper
}

var cache = cachehelper.New(CACHE.bb809_converter)

func (ch *CarHelper) GetCarInfoByCarID(carID string) (map[string]interface{}, error) {
	carCacheKey := RedisKeyHelper.GetCarCacheKey(carID)
	car, err := cache.Get(carCacheKey, CACHE.CAR_INFO)
	if err != nil {
		return nil, err
	}
	if car != nil && car["cnum"] != nil && car["plate_color"] != nil {
		return car, nil
	}
	car, err = CarRepo.GetCarByCarID(carID)
	if err != nil {
		return nil, err
	}
	if car == nil {
		return nil, nil
	}
	err = cache.Update(carCacheKey, car, CACHE.CAR_INFO)
	if err != nil {
		return nil, err
	}
	return car, nil
}

func (ch *CarHelper) SetFuelCutByCNum(cnum string, optType string, payload interface{}) error {
	return CarRepo.UpdateFuelCutByCNum(cnum, optType, payload)
}

func (ch *CarHelper) GetSettingsByCNum(cnum string) (map[string]interface{}, error) {
	return CarRepo.GetSettingsByCNum(cnum)
}
