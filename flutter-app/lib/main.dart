import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_blue/flutter_blue.dart';

void main() => runApp(MyApp());

FlutterBlue flutterBlue = FlutterBlue.instance;

final serviceGuid = Guid("00010000-0011-1000-8000-00805F9B34FB");
final chars = {
  "play": Guid("00010000-0011-1000-8000-00805F9B34FB"),
  "music": Guid("00060000-0011-1000-8000-00805F9B34FB"),
};


class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        // This is the theme of your application.
        //
        // Try running your application with "flutter run". You'll see the
        // application has a blue toolbar. Then, without quitting the app, try
        // changing the primarySwatch below to Colors.green and then invoke
        // "hot reload" (press "r" in the console where you ran "flutter run",
        // or simply save your changes to "hot reload" in a Flutter IDE).
        // Notice that the counter didn't reset back to zero; the application
        // is not restarted.
        primarySwatch: Colors.blue,
      ),
      home: MyHomePage(title: 'Flutter Demo Home Page'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  MyHomePage({Key key, this.title}) : super(key: key);

  // This widget is the home page of your application. It is stateful, meaning
  // that it has a State object (defined below) that contains fields that affect
  // how it looks.

  // This class is the configuration for the state. It holds the values (in this
  // case the title) provided by the parent (in this case the App widget) and
  // used by the build method of the State. Fields in a Widget subclass are
  // always marked "final".

  final String title;

  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class Music {
  Music(String n, [int d = 0, int f = 0]) : name = n, folder = d, file = f;
  String name = "";
  int folder;
  int file;
}

class _MyHomePageState extends State<MyHomePage> {
  int _status = 0;
  StreamSubscription<ScanResult> _scanSubscription;
  BluetoothDevice _device;
  BluetoothService _service;
  var _songs = <Music>[
    Music("haha"),
    Music("uh oh")
  ];

  BluetoothCharacteristic getChar(BluetoothService service, String charName) {
    return service?.characteristics.firstWhere((var c) {
      return c.uuid == chars[charName];
    });
  }

  Future write(String charName, List<int> data) async {
    var c = getChar(_service, charName);
    if (c == null) return null;
    return await c.write(data);
  }
  Future<List<int>> read(String charName) async {
    var c = getChar(_service, charName);
    if (c == null) return null;
    return await c.read();
  }


  void scanChars() async {
    if (_device == null || _scanSubscription == null) {
      return;
    }
    _scanSubscription.cancel();
    print('connecting...');
    await _device.connect();

    var _services = await _device.discoverServices();
    _services.forEach((var s) {
      if (s.uuid == serviceGuid) _service = s;
    });

    List<int> a = await read("music");
    _songs = [];
    a.asMap().forEach((var index, var value) {
      if (value == 255) {
        return;
      }
      for (int i = 1; i <= value; i++) {
        _songs.add(Music("${(index+1).toString().padLeft(2, '0')}/${i.toString().padLeft(3, '0')}.mp3", index+1, i));
      }
    });
    setState(() {
      _status = 2;
    });
  }
  void _onButtonPressed() async {
    if (_status == 1) {
      _scanSubscription.cancel();
      return;
    }

    _scanSubscription = flutterBlue.scan(withServices: [serviceGuid]).listen((scanResult) {
      // do something with scan result
      var device = scanResult.device;
      print('${device.name} found! rssi: ${scanResult.rssi}');
      _device = device;
      scanChars();
    });

    setState(() {
      _status = 1;
    });
  }

  void play(int folder, int file) async {
    await write("play", ['f'.codeUnitAt(0), folder, file]);
  }

  @override
  Widget build(BuildContext context) {
    // This method is rerun every time setState is called, for instance as done
    // by the _incrementCounter method above.
    //
    // The Flutter framework has been optimized to make rerunning build methods
    // fast, so that you can just rebuild anything that needs updating rather
    // than having to individually change instances of widgets.



    var items = <Widget>[];
    _songs.forEach((var music) {
//      items.add(
//        Padding(
//          padding: EdgeInsets.only(top: 16.0, right: 16.0, bottom: 16.0, left: 16.0),
//          child: Row(
//            children: <Widget>[
//                FlutterLogo(),
//                Text("eyy"),
//            ]
//          )
//        )
//      );
      items.add(
          Card(
            child: ListTile(
              leading: IconButton(
                  icon: Icon(Icons.play_circle_outline),
                  onPressed: () {
                    play(music.folder, music.file);
                  }
              ),
              title: Text(music.name)
            )
          )
      );
    });

    return Scaffold(
      body:
        Stack(
          children: <Widget>[
            Image(
              image: NetworkImage('https://www.billboard.com/files/styles/1500x992_gallery/public/media/drake-2019-bbyx-billboard-1548.jpg'),
            ),
            LayoutBuilder(
              builder: (BuildContext context, BoxConstraints viewportConstraints) {
                return SingleChildScrollView(
                    child: ConstrainedBox(
                        constraints: BoxConstraints(
                          minHeight: viewportConstraints.maxHeight,
                          minWidth: viewportConstraints.maxWidth,
                        ),
                        child: Padding(
                            padding: EdgeInsets.only(top: 200.0),
                            child: Container(
                                decoration: new BoxDecoration(
                                    color: Colors.white,
                                    borderRadius: new BorderRadius.only(
                                        topLeft: const Radius.circular(40.0),
                                        topRight: const Radius.circular(40.0)
                                    )
                                ),
                                child: Padding(
                                  padding: EdgeInsets.only(top: 32.0, right: 16.0, bottom: 16.0, left: 16.0),
                                  child: _status == 1 ? Center(child: CircularProgressIndicator()) : Column(
                                    children: items
                                  )
                                )
                            )
                        )
                    )
                );
                }
              ),
            ]

          ),
      floatingActionButton: FloatingActionButton(
        onPressed: _onButtonPressed,
        tooltip: 'Increment',
        child: Icon(_status == 0 ? Icons.bluetooth : (_status == 1 ? Icons.bluetooth_searching : Icons.bluetooth_connected)),
      ), // This trailing comma makes auto-formatting nicer for build methods.
    );
  }
}
