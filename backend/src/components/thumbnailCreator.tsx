import React, { useState } from 'react';
import ReactSlider from 'react-slider';
import styled from 'styled-components';

import { VideoToFrames, VideoToFramesMethod } from '../utils/VideoToFrame.js';

const StyledSlider = styled(ReactSlider)`
  width: 100%;
  height: 25px;
`;

const StyledThumb = styled.div`
  height: 25px;
  line-height: 25px;
  width: 25px;
  text-align: center;
  background-color: #000;
  color: #fff;
  border-radius: 50%;
  cursor: grab;
`;

function Thumb(props, state) {
  // eslint-disable-next-line react/destructuring-assignment, react/jsx-props-no-spreading
  return <StyledThumb {...props}>{state.valueNow}</StyledThumb>;
}

const StyledTrack = styled.div`
  top: 0;
  bottom: 0;
  background: #0f0;
  border-radius: 999px;
`;

function Track(props, state) {
  // eslint-disable-next-line react/jsx-props-no-spreading, react/destructuring-assignment
  return <StyledTrack {...props} index={state.index} />;
}

export default function ThumbnailCreator() {
  const [images, setImages] = useState([]);
  const [thumb, setThumb] = useState(0);
  const [status, setStatus] = useState('IDLE');
  const [frameCount, setFrameCount] = useState(300);

  const onInput = async (event) => {
    setImages([]);
    setStatus('LOADING');

    const [file] = event.target.files;
    const fileUrl = URL.createObjectURL(file);
    const frames = await VideoToFrames.getFrames(fileUrl, frameCount, VideoToFramesMethod.totalFrames);

    setStatus('IDLE');
    setImages(frames);
  };

  return (
    <div style={{
      width: '700px', border: '1px solid black', padding: '1em', borderRadius: '6px',
    }}
    >
      <h1>Thumbnail Creator</h1>
      <div style={{ padding: '1em' }}>
        {status === 'IDLE' ? 'Choose file: ' : 'Loading... '}
        <input type="file" className="visually-hidden" accept="video/*" onChange={onInput} />

        <div>
          Frame count:
          <input
            value={frameCount}
            onChange={(event) => {
              try {
                const k = parseInt(event.target.value, 10);
                setFrameCount(k);
              } catch (e) {
                console.log(`Error: ${e}`);
              }
            }}
          />
        </div>
      </div>

      {images?.length > 0 && (
        <div style={{ padding: '1em' }}>
          <StyledSlider
            defaultValue={0}
            min={0}
            max={images.length}
            renderTrack={Track}
            renderThumb={Thumb}
            onChange={(newValue, ThumbIndex) => {
              setThumb(newValue);
            }}
          />
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
            <img src={images[thumb]} alt="" />
          </div>
        </div>
      )}
    </div>
  );
}
