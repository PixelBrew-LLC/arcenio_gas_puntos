import 'dart:async';
import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:flutter_pos_printer_platform_image_3/flutter_pos_printer_platform_image_3.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// Información de una impresora Bluetooth descubierta o guardada.
class SavedPrinter {
  final String name;
  final String address;
  final bool isBle;

  SavedPrinter({required this.name, required this.address, this.isBle = false});
}

/// Singleton que gestiona la conexión Bluetooth con la impresora térmica.
/// La configuración se persiste en SharedPreferences y no se borra al cerrar sesión.
class PrinterService extends ChangeNotifier {
  static const _prefKeyName = 'printer_name';
  static const _prefKeyAddress = 'printer_address';
  static const _prefKeyIsBle = 'printer_is_ble';

  final _printerManager = PrinterManager.instance;

  // Estado
  bool _isScanning = false;
  bool _isConnected = false;
  SavedPrinter? _savedPrinter;
  List<PrinterDevice> _discoveredDevices = [];
  StreamSubscription<PrinterDevice>? _scanSubscription;
  StreamSubscription<BTStatus>? _btStatusSubscription;
  List<int>? _pendingTask;

  bool get isScanning => _isScanning;
  bool get isConnected => _isConnected;
  SavedPrinter? get savedPrinter => _savedPrinter;
  List<PrinterDevice> get discoveredDevices => _discoveredDevices;

  /// Inicializa el servicio: carga impresora guardada y escucha estado BT.
  Future<void> init() async {
    await _loadSavedPrinter();
    _listenBtStatus();
    // Intentar conectar automáticamente si hay impresora guardada
    if (_savedPrinter != null) {
      await connectToSaved();
    }
  }

  Future<void> _loadSavedPrinter() async {
    final prefs = await SharedPreferences.getInstance();
    final name = prefs.getString(_prefKeyName);
    final address = prefs.getString(_prefKeyAddress);
    final isBle = prefs.getBool(_prefKeyIsBle) ?? false;
    if (name != null && address != null) {
      _savedPrinter = SavedPrinter(name: name, address: address, isBle: isBle);
      notifyListeners();
    }
  }

  Future<void> _savePrinter(SavedPrinter printer) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_prefKeyName, printer.name);
    await prefs.setString(_prefKeyAddress, printer.address);
    await prefs.setBool(_prefKeyIsBle, printer.isBle);
    _savedPrinter = printer;
    notifyListeners();
  }

  Future<void> forgetPrinter() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_prefKeyName);
    await prefs.remove(_prefKeyAddress);
    await prefs.remove(_prefKeyIsBle);
    try {
      await _printerManager.disconnect(type: PrinterType.bluetooth);
    } catch (_) {}
    _savedPrinter = null;
    _isConnected = false;
    notifyListeners();
  }

  void _listenBtStatus() {
    _btStatusSubscription = _printerManager.stateBluetooth.listen((status) {
      if (status == BTStatus.connected) {
        _isConnected = true;
        // Si hay bytes pendientes, enviarlos
        if (_pendingTask != null) {
          final bytes = _pendingTask!;
          _pendingTask = null;
          final delay = Platform.isAndroid ? 1500 : 500;
          Future.delayed(Duration(milliseconds: delay), () {
            _printerManager.send(type: PrinterType.bluetooth, bytes: bytes);
          });
        }
      } else if (status == BTStatus.none) {
        _isConnected = false;
      }
      notifyListeners();
    });
  }

  /// Escanear impresoras Bluetooth.
  void scan() {
    _isScanning = true;
    _discoveredDevices = [];
    notifyListeners();

    _scanSubscription?.cancel();
    _scanSubscription = _printerManager
        .discovery(type: PrinterType.bluetooth, isBle: false)
        .listen(
          (device) {
            // Evitar duplicados
            if (!_discoveredDevices.any((d) => d.address == device.address)) {
              _discoveredDevices.add(device);
              notifyListeners();
            }
          },
          onDone: () {
            _isScanning = false;
            notifyListeners();
          },
          onError: (_) {
            _isScanning = false;
            notifyListeners();
          },
        );

    // Timeout de 10 segundos
    Future.delayed(const Duration(seconds: 10), () {
      if (_isScanning) {
        _scanSubscription?.cancel();
        _isScanning = false;
        notifyListeners();
      }
    });
  }

  /// Seleccionar y conectar a un dispositivo, guardarlo en preferences.
  Future<void> selectAndConnect(PrinterDevice device) async {
    // Desconectar si hay otra impresora conectada
    if (_isConnected) {
      try {
        await _printerManager.disconnect(type: PrinterType.bluetooth);
      } catch (_) {}
      _isConnected = false;
    }

    final printer = SavedPrinter(
      name: device.name,
      address: device.address!,
      isBle: false,
    );
    await _savePrinter(printer);
    await _connectTo(printer);
  }

  /// Conectar a la impresora guardada.
  Future<void> connectToSaved() async {
    if (_savedPrinter == null) return;
    await _connectTo(_savedPrinter!);
  }

  Future<void> _connectTo(SavedPrinter printer) async {
    try {
      await _printerManager.connect(
        type: PrinterType.bluetooth,
        model: BluetoothPrinterInput(
          name: printer.name,
          address: printer.address,
          isBle: printer.isBle,
          autoConnect: true,
        ),
      );
    } catch (e) {
      debugPrint('Error connecting to printer: $e');
    }
  }

  /// Enviar bytes ESC/POS a la impresora.
  Future<bool> printBytes(List<int> bytes) async {
    if (_savedPrinter == null) return false;

    if (_isConnected) {
      try {
        _printerManager.send(type: PrinterType.bluetooth, bytes: bytes);
        return true;
      } catch (e) {
        debugPrint('Print error: $e');
        return false;
      }
    } else {
      // Guardar bytes como tarea pendiente e intentar conectar
      _pendingTask = bytes;
      await connectToSaved();
      return true; // Se enviará cuando se conecte
    }
  }

  @override
  void dispose() {
    _scanSubscription?.cancel();
    _btStatusSubscription?.cancel();
    super.dispose();
  }
}
