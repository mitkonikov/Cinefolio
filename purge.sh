echo "Purging temporary files..."

# Backend purge
rm -rf ./backend/node_modules
rm -rf ./backend/dist
rm -rf ./backend/.adminjs

# Frontend purge
rm -rf ./cinefolio2/node_modules
rm -rf ./cinefolio2/.svelte-kit
rm -rf ./cinefolio2/build

rm -rf ./cinefolio/node_modules

echo "Purged all temporary files."
