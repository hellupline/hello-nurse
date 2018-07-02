for TAG in "landscape" "moon" "night" "scenic" "sky" "star" "sunset" "ruins"; do
    http POST http://query.hellupline.com/v1/tasks/nurse-fetch domain=konachan.net tag="${TAG}"
done
