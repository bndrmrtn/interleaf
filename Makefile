ATLAS_IMAGE := interleaf-atlas-pgvector:pg17
INPUT_FILE := testdata/input.pdf
OUTPUT_FILE := out/output.json

.PHONY: atlas-install atlas-diff atlas-apply atlas-image atlas-rebuild build run

atlas-rebuild:
	docker build -t $(ATLAS_IMAGE) ./docker/atlas

atlas-image:
	@docker image inspect $(ATLAS_IMAGE) >/dev/null 2>&1 || \
		$(MAKE) atlas-rebuild

atlas-install:
	@curl -sSf https://atlasgo.sh | sh

atlas-diff: atlas-image
	@atlas migrate diff --env gorm $(NAME)

build:
	@go build -o bin/ileaf ./cmd/ileaf

run: build
	@./bin/ileaf ${INPUT_FILE} ${OUTPUT_FILE}
ifdef SHOW
	@cat ${OUTPUT_FILE} | jq
endif
