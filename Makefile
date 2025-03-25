BIN = build/whitepointinventory
STAGING_BIN = build/whitepointinventorystaging

TMP_BIN = build/whitepointinventory.tmp
TMP_STAGING_BIN = build/whitepointinventorystaging.tmp

RED=\033[31m
GREEN=\033[32m
YELLOW=\033[33m
CYAN=\033[36m
RESET=\033[0m

.PHONY: build run clean

build:
	@echo "$(CYAN)Building production binary...$(RESET)"
	rm -f $(TMP_BIN)
	go build -ldflags "-X main.mountPath=/whitepoint" -o $(TMP_BIN) ./cmd
	mv $(TMP_BIN) $(BIN)
	@echo "$(GREEN)Build production binary successful$(RESET)"

build-staging:
	@echo "$(CYAN)Building staging binary...$(RESET)"
	rm -f $(TMP_STAGING_BIN)
	go build -ldflags "-X main.mountPath=/whitepointstaging" -o $(TMP_STAGING_BIN) ./cmd
	mv $(TMP_STAGING_BIN) $(STAGING_BIN)
	@echo "$(GREEN)Build staging binary successful$(RESET)"

run: clean build
	@echo "$(CYAN)Running production application locally...$(RESET)"
	$(BIN)

run-staging: clean-staging build-staging
	@echo "$(CYAN)Running staging application locally...$(RESET)"
	$(STAGING_BIN)

clean:
	@echo "$(CYAN)Cleaning production artifacts$(RESET)"
	rm -f $(BIN) $(TMP_BIN)
	@echo "$(YELLOW)Cleaned up production artifacts$(RESET)"

clean-staging:
	@echo "$(CYAN)Cleaning staging artifacts$(RESET)"
	rm -f $(STAGING_BIN) $(TMP_STAGING_BIN)
	@echo "$(YELLOW)Cleaned up staging artifacts$(RESET)"
