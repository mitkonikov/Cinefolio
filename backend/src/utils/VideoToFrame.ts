export enum VideoToFramesMethod {
  fps,
  totalFrames,
}

const WIDTH = 640;
const HEIGHT = 360;

export class VideoToFrames {
  /**
   * Extracts frames from the video and returns them as an array of imageData
   * @param videoUrl url to the video file (html5 compatible format) eg: mp4
   * @param amount number of frames per second or total number of frames that you want to extract
   * @param type [fps, totalFrames] The method of extracting frames:
   *             Number of frames per second of video or the total number of frames
   *             across the whole video duration. defaults to fps.
   */
  public static getFrames(
    videoUrl: string,
    amount: number,
    type: VideoToFramesMethod = VideoToFramesMethod.fps,
  ): Promise<string[]> {
    return new Promise((resolve: (frames: string[]) => void, reject: (error: string) => void) => {
      const frames: string[] = [];
      const canvas: HTMLCanvasElement = document.createElement('canvas');
      const context: CanvasRenderingContext2D = canvas.getContext('2d');
      let duration: number;

      const video = document.createElement('video');
      video.preload = 'auto';
      // eslint-disable-next-line @typescript-eslint/no-this-alias
      const that = this;
      video.addEventListener('loadeddata', async () => {
        canvas.width = WIDTH;
        canvas.height = HEIGHT;
        duration = video.duration;

        let totalFrames: number = amount;
        if (type === VideoToFramesMethod.fps) {
          totalFrames = duration * amount;
        }
        for (let time = 0; time < duration; time += duration / totalFrames) {
          // eslint-disable-next-line no-await-in-loop
          frames.push(await that.getVideoFrame(video, context, canvas, time));
        }
        resolve(frames);
      });
      video.src = videoUrl;
      video.load();
    });
  }

  private static getVideoFrame(
    video: HTMLVideoElement,
    context: CanvasRenderingContext2D,
    canvas: HTMLCanvasElement,
    time: number,
  ): Promise<string> {
    return new Promise((resolve: (frame: string) => void, reject: (error: string) => void) => {
      const eventCallback = () => {
        video.removeEventListener('seeked', eventCallback);
        this.storeFrame(video, context, canvas, resolve);
      };
      video.addEventListener('seeked', eventCallback);
      // eslint-disable-next-line no-param-reassign
      video.currentTime = time;
    });
  }

  private static storeFrame(
    video: HTMLVideoElement,
    context: CanvasRenderingContext2D,
    canvas: HTMLCanvasElement,
    resolve: (frame: string) => void,
  ) {
    context.drawImage(video, 0, 0, WIDTH, HEIGHT);
    resolve(canvas.toDataURL('image/jpeg', 0.5));
  }
}
