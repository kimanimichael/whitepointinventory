BIN = build/whitepointinventory
TMP_BIN = build/whitepointinventory.tmp

RED=\033[31m
GREEN=\033[32m
YELLOW=\033[33m
CYAN=\033[36m
RESET=\033[0m

.PHONY: build run clean

build:
	@echo "$(CYAN)Building application...$(RESET)"
	rm -f $(TMP_BIN)
	go build -o $(TMP_BIN) ./cmd
	mv $(TMP_BIN) $(BIN)
	@echo "$(GREEN)Build successful$(RESET)"

run: clean build
	@echo "$(CYAN)Running application...$(RESET)"
	$(BIN)

clean:
	@echo "$(CYAN)Cleaning artifacts$(RESET)"
	rm -f $(BIN) $(TMP_BIN)
	@echo "$(YELLOW)Cleaned up artifacts$(RESET)"
