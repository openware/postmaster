VERSION=`cat VERSION || echo 0.0.0`

LIST=(`echo $VERSION | tr '.' ' '`)
MAJOR=${LIST[0]}
MINOR=${LIST[1]}
PATCH=${LIST[2]}

PATCH=$((PATCH + 1))

TAG=$MAJOR.$MINOR.$PATCH
echo $TAG > VERSION

git add VERSION
git commit -m "v${TAG} Bump [ci skip]"
