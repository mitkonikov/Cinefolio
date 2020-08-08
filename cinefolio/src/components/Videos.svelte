<script>
    import { db } from './../js/firebase';

    const videos = db.collection('videos').orderBy('timestamp');
    console.log(videos);
</script>