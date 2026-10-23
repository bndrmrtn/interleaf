ATLAS_IMAGE := interleaf-atlas-pgvector:pg17

.PHONY: atlas-install atlas-diff atlas-apply atlas-image atlas-rebuild

atlas-rebuild:
	docker build -t $(ATLAS_IMAGE) ./docker/atlas

atlas-image:
	@docker image inspect $(ATLAS_IMAGE) >/dev/null 2>&1 || \
		$(MAKE) atlas-rebuild

atlas-install:
	@curl -sSf https://atlasgo.sh | sh

atlas-diff: atlas-image
	@atlas migrate diff --env gorm $(NAME)
