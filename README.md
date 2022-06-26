# EmeraldCity-CLI

:warning: THIS PROJECT IS UNDER ACTIVE DEVELOPMENT AND NOT READY FOR PRODUCTION USAGE!

The EmeraldCity-CLI project was developed to allow the [Emerald City Playground](https://github.com/BoiseITGuru/tx-script-utility/tree/dev-local-emulator) to start and manage a local Flow Emulator. "Phase 2" of the project will utlize the existing WebSocket connection to allow developers to utilize their local workstation as a file system for the Emerald City Playground.

## Challenges

Originally the plan was to utilize [Overflow](https://github.com/bjartek/overflow) to start and manage a Flow Emulator. However, as development procedeed and the other uses cases of the CLI started becoming clear it was decided to directly to import a slightly customized version of the [Flow Emulator](https://github.com/BoiseITGuru/flow-emulator) that brings the [graceland](https://github.com/psiemens/graceland) multi-routine mangement forward into the startup arguments, also allowing for us to pass the [flow-dev-wallet](https://github.com/BoiseITGuru/fcl-dev-wallet/tree/dev-full-accounts) into the emulator startup routine. 

The flow-dev-wallet first needed to be modified to support graceland then the accounts functions had to be re-written to use actuall accounts. It turns out the dev-wallet does not create any addtional accounts, it just uses a special contacts to create alias accounts for the service account. All functions in ```src/accounts.ts``` of the dev-wallet are being re-written to use an API running alongside the dev-wallet, the API will utlize Overflow/flow-do-sdk to provide a more realistic dev-wallet

## Current Status

The playground and CLI are still undergoing heavy development, the CLI is written in GO and requires manual compiling still. Once we have finished development it will be available as an NPM package. Currently, the CLI connects to the playground and starts the emulators as well as executes scripts.

The dev-wallet modification should be done shortly once that is completed we can natively use FCL for all interactions with the playground OR **any other project!** the changes to the emulator and dev-wallet are fully backward compatible with the Flow CLI allowing for use in any project with a single command.

## How to use the EmeraldCity-CLI

Current Commands **propose we create the shortname "ecd" to invoke the CLI **

1. ```EmeraldCity-CLI emulator```
    * Starts a WebSocket server to await connections from the Emerald City Playground and upon connection start the Flow Emulator and Dev Wallet. The defualt confiugration starts the Flow Emulator with the ```--contracts``` flag set to automatically deploy all base contracts.
        * Flags - all flags for the Flow Emulator are supported with the following exceptions/addtions:
            1. The ```--contracts``` flag has been renamed to ```--no-contracts``` and will start the emulator without the NonFungibleToken, MetadataViews, ExampleNFT, or NFTStorefront contracts.
            2. ```--no-ws-server``` will start the Flow Emulator without creating a WebSocket server to wait for the Emerald City Playground to connect.
            3. ```--no-dev-wallet``` will start the Flow Emulator without start the Flow Dev Wallet
            4. ```--file-system, -f``` will connect the current directory as a file system to the Emerald City Playground **COMING SOON - IN DEVELOPMENT**
    1. ```EmeraldCity-CLI emulator listAccounts```
        * List all accounts for the currently running emulator
2. ```EmeraldCity-CLI create-fcl-app <nextjs/svelte>``` **COMING SOON - IN DEVELOPMENT**
    * Deploys FCL Quick Start project code sets chossing either the fcl-nextjs-quickstart or the fcl-svelte-quickstart project templates
        1. ```--name <name>```