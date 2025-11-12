package loading

import "e-voting-service/data/configuration"

func generalizedLoaderFactory[T any](db_loader T, mock_loader T) T {
	conf_mock := configuration.GlobalConfig.Use_mock_data
	if conf_mock {
		return mock_loader
	}
	return db_loader
}
