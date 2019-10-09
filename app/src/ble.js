
function getUUID(n) {
	return n.toString(16).padStart(8, "0").slice(-8) + "-0011-1000-8000-00805F9B34FB".toLowerCase()
}

const SERVICES = [
	getUUID(0x00010000)
];
const CHARS = {
  TestSvcUUID    : 0x00010000,
	PlayCharUUID   : 0x00020000,
	VolumeCharUUID : 0x00030000,
	LoopCharUUID   : 0x00040000,
	EQCharUUID     : 0x00050000,
	MusicCharUUID  : 0x00060000,
	StatusCharUUID : 0x00070000,
	TestCharUUID   : 0x00080000
}


function getSupportedProperties(characteristic) {
  let supportedProperties = [];
  for (const p in characteristic.properties) {
    if (characteristic.properties[p] === true) {
      supportedProperties.push(p.toUpperCase());
    }
  }
  return '[' + supportedProperties.join(', ') + ']';
}

async function getServer() {
  // Validate services UUID entered by user first.
  /*let optionalServices = document.querySelector('#optionalServices').value
    .split(/, ?/).map(s => s.startsWith('0x') ? parseInt(s) : s)
    .filter(s => s && BluetoothUUID.getService);
  */
  console.log('Requesting any Bluetooth Device...');
  let device = await navigator.bluetooth.requestDevice({
    filters: [{
      services: SERVICES
    }] // <- Prefer filters to save energy & show relevant devices.
    //acceptAllDevices: true,
    //optionalServices: optionalServices
  })
  console.log('Connecting to GATT Server...');
  console.log("gatt:", device.name)
  return device.gatt.connect();
}
function getService(server) {
  console.log('Getting Services...');
  return server.getPrimaryService(SERVICES[0])
}
function getChars(service) {
  return service.getCharacteristics()
}



class BLE {
  static SERVICES = SERVICES
  static CHARS = CHARS
  server = null
  service = null
  chars = []
  isConnected() {
    return this.server && this.service && this.chars.length > 0
  }
  async reconnect() {
    this.server = null;
    this.service = null;
    this.chars = [];
    this.fetchChars();
  }
  async fetchChars() {
    this.server = this.server || await getServer();
    this.service = this.service || await getService(this.server);
    this.chars = await getChars(this.service);
    
    console.log('> Service: ' + this.service.uuid);
    this.chars.forEach(characteristic => {
      console.log('>> Characteristic: ' + characteristic.uuid + ' ' + getSupportedProperties(characteristic));
    });
    return true;
  }
  getChar(n) {
    return this.chars.find(v => v.uuid === getUUID(n).toLowerCase());
  }
}


export default BLE

