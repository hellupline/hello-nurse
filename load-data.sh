for TAG in "landscape" "moon" "night" "scenic" "sky" "star" "sunset" "ruins"; do
    http POST :8080/v1/tasks/booru/fetch-tag domain=konachan.net tag="${TAG}"
done
