import React, {useState} from 'react';
import { Switch, Route, Link } from "react-router-dom";
import { FiChevronLeft, FiChevronRight, FiPlay, FiBluetooth, FiVolume2 } from "react-icons/fi";

import BLE from './ble'
import logo from './logo.svg';
import './App.css';

const ble = new BLE();

function Header(props) {
  return <header className="header">
    <div style={{fontSize: "1.5em"}}>Started from the bottom now we're here</div>
  </header>
}
function Player() {
  const [progress, setProgress] = useState(0);
  return <div className="player">
    <div style={{position: "relative", top: 0, height: "4px", width: "100%", overflow: "hidden"}}>
      <div className="player-progressBar" style={{
        transform: `translateX(${progress}%)`
      }}>
      </div>
    </div>
    <div style={{padding: "1em 1em", display: "flex", alignItems: "center"}}>
      <FiChevronLeft style={{marginLeft: "1em"}}/>
      <FiPlay style={{marginLeft: "1em"}}/>
      <FiChevronRight style={{marginLeft: "1em"}}/>
      <div style={{marginLeft: "auto", overflow: "hidden"}}>
        sogn name
      </div>
      <FiVolume2 style={{marginLeft: "auto", marginRight: "1em"}}/>
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

  function connect() {
    ble.connect().then(() => {
      setConnecting(true)
      ble.chars.forEach(characteristic => {console.log(characteristic)});
      ble.getChar(BLE.CHARS.VolumeCharUUID).readValue().then(r => {
        console.log(r);
      })
      ble.getChar(BLE.CHARS.MusicCharUUID).readValue().then(r => {
        console.log(r);
      })
    }).catch(e => {
      console.log(e)
    })
  }

  return (
    <div className="App">
      <Header/>
      <FiBluetooth style={{marginTop: "0.5em", color: connecting ? "#1976d2" : "white"}}/>
      <div className="buttonContainer">
        <div style={{position: "relative", width: "100%", display: "flex", flexFlow: "column", justifyContent: "center", alignItems: "center"}}>
          <div className="startButton" onClick={() => connect()}>
          </div>
          <div style={{height: list && list.length > 0 ? 'auto' : 0, transition: "height 1s"}}>
            {
              list.map((v, i) => {
                return <div key={i} className="item">
                  {v.uuid}
                </div>
              })
            }
          </div>
        </div>
      </div>
      <Player/>
    </div>
  );
}

export default App;
