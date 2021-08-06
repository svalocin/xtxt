mkdir _push
cd _push
git clone https://x-access-token:$1@github.com/slyerr/xtxt.git
# git clone https://github.com/slyerr/xtxt.git

cd xtxt
git config user.email "actions@github.com"
git config user.name "Actions"
git checkout -b publish origin/publish
cp -rf ../../out ../_xtxt
cp -rf .git ../_xtxt

cd ../_xtxt
if [ -n "$(git status -s)" ]; then
  git add -f .
  git commit -m "chore: Actions update 喷嚏网 RSS"
  git push
fi
