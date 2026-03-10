import 'package:esc_pos_utils/esc_pos_utils.dart';
import 'package:flutter/services.dart';
import 'package:image/image.dart' as img;
import 'package:intl/intl.dart';
import '../utils/cedula_utils.dart';

/// Genera bytes ESC/POS para tickets de acumulación y canje.
class TicketGenerator {
  static const _paperSize = PaperSize.mm58;
  static const _profileName = 'XP-N160I';

  /// Genera ticket de ACUMULACIÓN de puntos.
  static Future<List<int>> generateEarnTicket({
    required String clientName,
    required String clientCedula,
    required double gallons,
    required double pointsEarned,
    required double newBalance,
    required DateTime date,
  }) async {
    final profile = await CapabilityProfile.load(name: _profileName);
    final generator = Generator(_paperSize, profile);
    List<int> bytes = [];

    bytes += generator.setGlobalCodeTable('CP1252');

    // Logo
    bytes += await _addLogo(generator);
    bytes += generator.feed(1);

    // Header
    bytes += _addHeader(generator);
    bytes += _addSeparator(generator);

    bytes += generator.text(
      'TICKET DE ACUMULACION',
      styles: const PosStyles(align: PosAlign.center, bold: true),
    );
    bytes += generator.text(
      'Fecha: ${DateFormat('dd/MM/yyyy hh:mm a').format(date)}',
      styles: const PosStyles(align: PosAlign.center),
    );

    bytes += _addSeparator(generator);
    bytes += generator.feed(1);

    // Datos del cliente
    bytes += generator.text(
      'DATOS DEL CLIENTE:',
      styles: const PosStyles(align: PosAlign.left, bold: true),
    );
    bytes += generator.text(
      'Cedula: ${CedulaUtils.format(clientCedula)}',
      styles: const PosStyles(align: PosAlign.left),
    );
    bytes += generator.text(
      'Nombre: $clientName',
      styles: const PosStyles(align: PosAlign.left),
    );

    bytes += generator.feed(1);
    bytes += _addSeparator(generator);

    // Detalles
    bytes += generator.text(
      'DETALLES DE LA COMPRA:',
      styles: const PosStyles(align: PosAlign.left, bold: true),
    );
    bytes += generator.text(
      'Galones: ${gallons.toStringAsFixed(2)} gal',
      styles: const PosStyles(align: PosAlign.left),
    );

    bytes += generator.feed(1);
    bytes += _addSeparator(generator);

    // Puntos ganados
    bytes += generator.text(
      'PUNTOS GANADOS:',
      styles: const PosStyles(align: PosAlign.center, bold: true),
    );
    bytes += generator.feed(1);
    bytes += generator.text(
      '${pointsEarned.toStringAsFixed(0)} PUNTOS',
      styles: const PosStyles(
        align: PosAlign.center,
        bold: true,
        height: PosTextSize.size2,
        width: PosTextSize.size1,
      ),
    );
    bytes += generator.feed(1);
    bytes += _addSeparator(generator);

    // Saldo
    bytes += generator.text(
      'SALDO ACUMULADO:',
      styles: const PosStyles(align: PosAlign.center, bold: true),
    );
    bytes += generator.feed(1);
    bytes += generator.text(
      '${newBalance.toStringAsFixed(0)} PUNTOS',
      styles: const PosStyles(
        align: PosAlign.center,
        bold: true,
        height: PosTextSize.size2,
        width: PosTextSize.size1,
      ),
    );

    bytes += generator.feed(1);
    bytes += _addSeparator(generator);

    bytes += _addFooter(generator);

    bytes += generator.feed(3);
    bytes += generator.cut();

    return bytes;
  }

  /// Genera ticket de CANJE de puntos.
  static Future<List<int>> generateRedeemTicket({
    required String clientName,
    required String clientCedula,
    required double pointsRedeemed,
    required double newBalance,
    required DateTime date,
  }) async {
    final profile = await CapabilityProfile.load(name: _profileName);
    final generator = Generator(_paperSize, profile);
    List<int> bytes = [];

    bytes += generator.setGlobalCodeTable('CP1252');

    // Logo
    bytes += await _addLogo(generator);
    bytes += generator.feed(1);

    // Header
    bytes += _addHeader(generator);
    bytes += _addSeparator(generator);

    bytes += generator.text(
      'TICKET DE CANJE',
      styles: const PosStyles(align: PosAlign.center, bold: true),
    );
    bytes += generator.text(
      'Fecha: ${DateFormat('dd/MM/yyyy hh:mm a').format(date)}',
      styles: const PosStyles(align: PosAlign.center),
    );

    bytes += _addSeparator(generator);
    bytes += generator.feed(1);

    // Datos del cliente
    bytes += generator.text(
      'DATOS DEL CLIENTE:',
      styles: const PosStyles(align: PosAlign.left, bold: true),
    );
    bytes += generator.text(
      'Cedula: ${CedulaUtils.format(clientCedula)}',
      styles: const PosStyles(align: PosAlign.left),
    );
    bytes += generator.text(
      'Nombre: $clientName',
      styles: const PosStyles(align: PosAlign.left),
    );

    bytes += generator.feed(1);
    bytes += _addSeparator(generator);

    // Puntos canjeados
    bytes += generator.text(
      'PUNTOS CANJEADOS:',
      styles: const PosStyles(align: PosAlign.center, bold: true),
    );
    bytes += generator.feed(1);
    bytes += generator.text(
      '${pointsRedeemed.toStringAsFixed(0)} PUNTOS',
      styles: const PosStyles(
        align: PosAlign.center,
        bold: true,
        height: PosTextSize.size2,
        width: PosTextSize.size1,
      ),
    );
    bytes += generator.feed(1);
    bytes += _addSeparator(generator);

    // Nuevo saldo
    bytes += generator.text(
      'NUEVO SALDO:',
      styles: const PosStyles(align: PosAlign.center, bold: true),
    );
    bytes += generator.feed(1);
    bytes += generator.text(
      '${newBalance.toStringAsFixed(0)} PUNTOS',
      styles: const PosStyles(
        align: PosAlign.center,
        bold: true,
        height: PosTextSize.size2,
        width: PosTextSize.size1,
      ),
    );

    bytes += generator.feed(1);
    bytes += _addSeparator(generator);

    bytes += _addFooter(generator);

    bytes += generator.feed(3);
    bytes += generator.cut();

    return bytes;
  }

  // --- Helpers ---

  static Future<List<int>> _addLogo(Generator generator) async {
    try {
      final ByteData data = await rootBundle.load(
        'assets/arcenio_logo_printer.png',
      );
      final Uint8List bytesImg = data.buffer.asUint8List();
      img.Image? image = img.decodeImage(bytesImg);
      if (image != null) {
        int newWidth = 350;
        int newHeight = ((image.height / image.width) * newWidth).round();
        image = img.copyResize(image, width: newWidth, height: newHeight);
        return generator.imageRaster(image, align: PosAlign.center);
      }
    } catch (_) {}
    return [];
  }

  static List<int> _addHeader(Generator generator) {
    List<int> bytes = [];
    bytes += generator.text(
      'Estacion de Servicio Arcenio',
      styles: const PosStyles(align: PosAlign.center, bold: true),
    );
    bytes += generator.text(
      'Calle Principal, Villa Riva',
      styles: const PosStyles(align: PosAlign.center),
    );
    bytes += generator.text(
      'Tel: (809) 234-5678',
      styles: const PosStyles(align: PosAlign.center),
    );
    return bytes;
  }

  static List<int> _addSeparator(Generator generator) {
    return generator.text(
      '--------------------------------',
      styles: const PosStyles(align: PosAlign.center),
    );
  }

  static List<int> _addFooter(Generator generator) {
    List<int> bytes = [];
    bytes += generator.text(
      'Gracias por su preferencia!',
      styles: const PosStyles(align: PosAlign.center, bold: true),
    );
    return bytes;
  }
}
