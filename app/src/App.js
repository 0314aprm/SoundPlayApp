import React, {useState} from 'react';
import { Switch, Route, Link } from "react-router-dom";
import { FiChevronLeft, FiChevronRight, FiPlay, FiPause, FiBluetooth, FiVolume1, FiVolume2, FiVolumeX, FiVolume } from "react-icons/fi";

import BLE from './ble'
import logo from './logo.svg';
import './App.css';

const ble = new BLE();

function getFileName(n, t = 2) {
  return n.toString().padStart(t, "0");
}
function getPath(folder, file) {
  return getFileName(folder) + "/" + getFileName(file, 3);
}
function Uint16ToArray(v) {
  return Uint8Array.of(v >> 8, v)
}

function Header(props) {
  return <header className="header">
    <div style={{fontSize: "1.5em"}}>Started from the bottom now we're here</div>
    <FiBluetooth style={{marginTop: "0.5em", color: props.connecting ? "#1976d2" : "white"}} onClick={() => {props.reconnect()}}/>
  </header>
}

function Volume(props) {
  if (props.volume === 0) {
    return <FiVolumeX style={{marginLeft: "auto", marginRight: "1em"}}/>
  } else if (props.volume <= 10) {
    return <FiVolume style={{marginLeft: "auto", marginRight: "1em"}}/>
  } else if (props.volume <= 20) {
    return <FiVolume1 style={{marginLeft: "auto", marginRight: "1em"}}/>
  } else if (props.volume <= 30) {
    return <FiVolume2 style={{marginLeft: "auto", marginRight: "1em"}}/>
  }
}
function Player(props) {
  const [progress, setProgress] = useState(0);

  return <div className="player">
    <div style={{position: "relative", top: 0, height: "4px", width: "100%", overflow: "hidden"}}>
      <div className="player-progressBar" style={{
        transform: `translateX(${progress}%)`
      }}>
      </div>
    </div>
    <div style={{padding: "1em 1em", display: "flex", alignItems: "center"}}>
      <FiChevronLeft className="player-control" onClick={props.previous}/>
      {
        props.playing ?
          <FiPause className="player-control" onClick={props.resume}/>
        :
          <FiPlay className="player-control" onClick={props.pause}/>
      }
      <FiChevronRight className="player-control" onClick={props.next}/>
      <div style={{marginLeft: "auto", overflow: "hidden"}}>
        {props.songName}
      </div>
      <div className="player-control" onClick={props.toggleVolume}>
        <Volume volume={props.volume}/>
      </div>
    </div>
  </div>
}
function Content() {
  return <main className="content">
    <Switch>
      <Route path="/">
        <Explore />
      </Route>
      <Route path="/status">
        <Status/>
      </Route>
    </Switch>
  </main>
}
function Explore() {
  const [songs, setSongs] = useState(["dwa", "ooo"])
  return <div>
    {
      songs.map((v, i) => {
        return <div key={i}>
          {v}
        </div>
      })
    }
  </div>
}
function Status() {
  return <div>
  </div>
}



function App() {
  const [list, setList] = useState([]);
  const [folders, setFolders] = useState([[]]);
  const [connecting, setConnecting] = useState(false);
  const [currentSong, setCurrentSong] = useState("");
  const [volume, setVolume] = useState(0);
  const [loop, setLoop] = useState("");
  const [playing, setPlaying] = useState(false);



  async function connect() {
    try {
      await ble.fetchChars();
      setConnecting(true)
      loadChars()
      fetchVolume()
      fetchMusics()
    } catch (e) {
      console.log(e)
    }
  }
  async function reconnect() {
    try {
      await ble.reconnect();
      setConnecting(true)
      loadChars()
      fetchVolume()
      fetchMusics()
    } catch (e) {
      console.log(e)
    }
  }

  function fetchVolume() {
    ble.getChar(BLE.CHARS.VolumeCharUUID).readValue().then(r => {
      console.log("volume:", r.getUint8(0));
      setVolume(r.getUint8(0));
    })
  }
  function fetchMusics() {
    console.log("fecth")
    console.log("fecth", ble.getChar(BLE.CHARS.MusicCharUUID))
    ble.getChar(BLE.CHARS.MusicCharUUID).readValue().then(r => {console.log(r)})
    console.log("fecth2")
    ble.getChar(BLE.CHARS.MusicCharUUID).readValue().then(r => {
      console.log("fetc,usic", r)
      let newFolder = [];
      for (let i = 0; i < r.byteLength; i++) {
        let files = r.getUint8(i);
        if (files === 0xFF) continue;
        
        newFolder.push(files);
        console.log("folder:", i + 1, " -> ", files);
      }
      setFolders(newFolder);
    }).catch(e => {
      console.log("music: ", e)
    })
  }
  async function loadChars() {
    ble.chars.forEach(characteristic => {console.log(characteristic)});
  }


  function sendPlay(data) {
    ble.getChar(BLE.CHARS.PlayCharUUID).writeValue(data).then(r => {
      console.log(r);
    })
  }
  function playMusic(folder, file) {
    console.log("playing ", folder, "/", file)
    sendPlay(Uint8Array.of('f'.charCodeAt(), folder, file))
    setCurrentSong(getPath(folder, file))
  }
  function resume() {
    console.log("resume")
    setPlaying(true)
    sendPlay(Uint8Array.of('r'.charCodeAt()))
  }
  function pause() {
    console.log("pause")
    setPlaying(false)
    sendPlay(Uint8Array.of('s'.charCodeAt()))
  }
  function next() {
    console.log("next")
    setPlaying(true)
    sendPlay(Uint8Array.of('n'.charCodeAt()))
  }
  function previous() {
    console.log("previous")
    setPlaying(true)
    sendPlay(Uint8Array.of('p'.charCodeAt()))
  }
  function sendVolume(volume) {
    console.log("set volume to", volume)
    setVolume(volume)
    ble.getChar(BLE.CHARS.VolumeCharUUID).writeValue(Uint16ToArray(volume)).then(r => {
      console.log(r);
    })
  }
  function toggleVolume() {
    // 0 - 10 - 20 - 30
    let v = (volume - volume % 10 + 10) % 40;
    console.log("toggleVolume: ", v)
    sendVolume(v)
  }

  return (
    <div className="App">
      <Header connecting={connecting} reconnect={reconnect}/>
      <div className="buttonContainer">
        <div style={{position: "relative", width: "100%", display: "flex", flexFlow: "column", justifyContent: "center", alignItems: "center"}}>
          <div className="startButton" onClick={() => connect()}>
          </div>
          <div style={{transition: "height 1s"}}>
            {
              folders.map((v, folder) => {
                return <div key={folder} className="item">
                  <div>
                    {getFileName(folder + 1)}
                  </div>
                  <div>
                    {
                      (() => {
                        let a = [];
                        for (let i = 0; i < v; i++) {
                          a.push(<div key={i} onClick={() => {playMusic(folder + 1, i + 1)}}>
                            {
                              getFileName(i + 1) + ".mp3"
                            }
                          </div>)
                        }
                        return a;
                      })()
                    }
                  </div>
                </div>
              })
            }
          </div>
        </div>
      </div>
      {
        <Player 
          songName={currentSong}
          volume={volume}
          playing={playing}
          toggleVolume={toggleVolume}
          next={next}
          previous={previous}
          resume={resume}
          pause={pause}
        />
      }
    </div>
  );
}

export default App;
