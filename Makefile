.PHONY: assets

assets:
	cd internal/server && statik -src=static
