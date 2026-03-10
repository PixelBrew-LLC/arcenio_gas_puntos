import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_mlkit_text_recognition/google_mlkit_text_recognition.dart';
import 'package:image_picker/image_picker.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:toastification/toastification.dart';
import '../../../core/network/dio_client.dart';
import '../../../core/printer/printer_service.dart';
import '../../../core/printer/ticket_generator.dart';
import '../../../core/theme.dart';
import '../../../core/utils/cedula_utils.dart';
import '../data/client_remote_datasource.dart';

class ClientSearchPage extends StatefulWidget {
  const ClientSearchPage({super.key});

  @override
  State<ClientSearchPage> createState() => _ClientSearchPageState();
}

class _ClientSearchPageState extends State<ClientSearchPage> {
  final _cedulaController = TextEditingController();
  bool _isLoading = false;
  bool _isScanningCamera = false;
  Map<String, dynamic>? _client;
  String? _errorMessage;
  double? _balance;
  double _minRedeem = 0;
  List<dynamic> _history = [];

  @override
  void dispose() {
    _cedulaController.dispose();
    super.dispose();
  }

  Future<void> _scanCedulaFromCamera() async {
    try {
      final picker = ImagePicker();
      final XFile? image = await picker.pickImage(source: ImageSource.camera);
      if (image == null) return;

      setState(() => _isScanningCamera = true);

      final inputImage = InputImage.fromFilePath(image.path);
      final recognizer = TextRecognizer(script: TextRecognitionScript.latin);
      final RecognizedText recognized = await recognizer.processImage(
        inputImage,
      );
      await recognizer.close();

      // Extract only the cedula number (format 000-0000000-0)
      final cedulaRegex = RegExp(r'\d{3}-\d{7}-\d');
      final match = cedulaRegex.firstMatch(recognized.text);

      if (match != null) {
        _cedulaController.text = match.group(0)!;
      } else {
        // Try without dashes: 11 digit number
        final rawRegex = RegExp(r'\b\d{11}\b');
        final rawMatch = rawRegex.firstMatch(
          recognized.text.replaceAll(' ', ''),
        );
        if (rawMatch != null) {
          _cedulaController.text = CedulaUtils.format(rawMatch.group(0)!);
        } else {
          if (mounted) {
            toastification.show(
              context: context,
              type: ToastificationType.warning,
              style: ToastificationStyle.flatColored,
              title: const Text('No se encontró cédula'),
              description: const Text(
                'No se pudo detectar un número de cédula en la imagen.',
              ),
              alignment: Alignment.topCenter,
              autoCloseDuration: const Duration(seconds: 4),
              showProgressBar: false,
            );
          }
        }
      }
    } catch (e) {
      if (mounted) {
        toastification.show(
          context: context,
          type: ToastificationType.error,
          style: ToastificationStyle.flatColored,
          title: const Text('Error al escanear'),
          description: const Text(
            'Hubo un error al leer la cédula desde la cámara.',
          ),
          alignment: Alignment.topCenter,
          autoCloseDuration: const Duration(seconds: 4),
          showProgressBar: false,
        );
      }
    } finally {
      if (mounted) setState(() => _isScanningCamera = false);
    }
  }

  Future<void> _search() async {
    final rawCedula = CedulaUtils.unformat(_cedulaController.text);
    if (rawCedula.isEmpty) return;

    final error = CedulaUtils.validate(rawCedula);
    if (error != null) {
      setState(() => _errorMessage = error);
      return;
    }

    setState(() {
      _isLoading = true;
      _errorMessage = null;
      _client = null;
      _balance = null;
      _minRedeem = 0;
      _history = [];
    });

    try {
      final dioClient = context.read<DioClient>();
      final datasource = ClientRemoteDatasource(dioClient);
      final result = await datasource.searchByCedula(rawCedula);

      // Obtener balance
      double balance = 0;
      double minRedeem = 0;
      try {
        final balanceResult = await dioClient.dio.get(
          '/transactions/balance/${result['id']}',
        );
        balance = (balanceResult.data['balance'] as num).toDouble();
        minRedeem = (balanceResult.data['min_redeem'] as num?)?.toDouble() ?? 0;
      } catch (_) {}

      // Obtener historial
      List<dynamic> history = [];
      try {
        final historyResult = await dioClient.dio.get(
          '/transactions/history/${result['id']}',
        );
        history = historyResult.data as List<dynamic>;
      } catch (_) {}

      if (mounted) {
        setState(() {
          _client = result;
          _balance = balance;
          _minRedeem = minRedeem;
          _history = history;
        });
      }
    } catch (e) {
      final message = e.toString().replaceAll('Exception: ', '');
      final isNotFound = message.toLowerCase().contains('no encontrado');

      if (mounted) {
        if (isNotFound) {
          setState(() => _isLoading = false);
          _showClientNotFoundDialog(rawCedula);
          return;
        } else {
          setState(() => _errorMessage = message);
        }
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Future<void> _showClientNotFoundDialog(String rawCedula) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Cliente no encontrado'),
        content: Text(
          'No existe un cliente con la cédula ${CedulaUtils.format(rawCedula)}.\n\n¿Desea registrarlo ahora?',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(false),
            child: const Text('Cancelar'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(ctx).pop(true),
            child: const Text('Registrar'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      context.push('/clients/register', extra: {'initialCedula': rawCedula});
    }
  }

  String _formatDate(String isoDate) {
    try {
      final date = DateTime.parse(isoDate).toLocal();
      return DateFormat('dd/MM/yyyy HH:mm').format(date);
    } catch (_) {
      return isoDate;
    }
  }

  String _formatExpiryDate(String isoDate) {
    try {
      final date = DateTime.parse(isoDate).toLocal();
      final now = DateTime.now();
      final isExpired = date.isBefore(now);
      final formatted = DateFormat('dd/MM/yyyy').format(date);
      return isExpired ? '$formatted (vencido)' : formatted;
    } catch (_) {
      return isoDate;
    }
  }

  bool _isExpired(String? isoDate) {
    if (isoDate == null) return false;
    try {
      return DateTime.parse(isoDate).toLocal().isBefore(DateTime.now());
    } catch (_) {
      return false;
    }
  }

  Future<void> _reprintTransaction(dynamic tx) async {
    if (_client == null) return;

    final type = tx['transaction_type'] as String;
    final isEarn = type == 'earn';
    final points = (tx['points'] as num).toDouble().abs();
    final clientName = '${_client!['nombres']} ${_client!['apellidos']}';
    final clientCedula = _client!['cedula'] as String;

    try {
      List<int> bytes;
      if (isEarn) {
        final gallons = (tx['gallons_amount'] as num?)?.toDouble() ?? 0;
        bytes = await TicketGenerator.generateEarnTicket(
          clientName: clientName,
          clientCedula: clientCedula,
          gallons: gallons,
          pointsEarned: points,
          newBalance: _balance ?? 0,
          date:
              DateTime.tryParse(tx['created_at'] ?? '')?.toLocal() ??
              DateTime.now(),
        );
      } else {
        bytes = await TicketGenerator.generateRedeemTicket(
          clientName: clientName,
          clientCedula: clientCedula,
          pointsRedeemed: points,
          newBalance: _balance ?? 0,
          date:
              DateTime.tryParse(tx['created_at'] ?? '')?.toLocal() ??
              DateTime.now(),
        );
      }

      if (mounted) {
        final printerService = context.read<PrinterService>();
        await printerService.printBytes(bytes);
        toastification.show(
          context: context,
          type: ToastificationType.success,
          style: ToastificationStyle.flatColored,
          title: const Text('Reimprimiendo ticket'),
          description: const Text('El ticket se está enviando a la impresora.'),
          alignment: Alignment.topCenter,
          autoCloseDuration: const Duration(seconds: 3),
          showProgressBar: false,
        );
      }
    } catch (_) {
      if (mounted) {
        toastification.show(
          context: context,
          type: ToastificationType.error,
          style: ToastificationStyle.flatColored,
          title: const Text('Error al reimprimir'),
          description: const Text(
            'No se pudo enviar el ticket a la impresora.',
          ),
          alignment: Alignment.topCenter,
          autoCloseDuration: const Duration(seconds: 4),
          showProgressBar: false,
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Buscar Cliente')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Search bar
            Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _cedulaController,
                    decoration: InputDecoration(
                      labelText: 'Cédula del cliente',
                      prefixIcon: const Icon(Icons.badge_outlined),
                      hintText: '###-#######-#',
                      suffixIcon: _isScanningCamera
                          ? const Padding(
                              padding: EdgeInsets.all(12),
                              child: SizedBox(
                                width: 20,
                                height: 20,
                                child: CircularProgressIndicator(
                                  strokeWidth: 2,
                                ),
                              ),
                            )
                          : IconButton(
                              icon: const Icon(Icons.camera_alt_outlined),
                              tooltip: 'Escanear cédula con cámara',
                              onPressed: _scanCedulaFromCamera,
                            ),
                    ),
                    keyboardType: TextInputType.number,
                    inputFormatters: [CedulaInputFormatter()],
                    textInputAction: TextInputAction.search,
                    onSubmitted: (_) => _search(),
                  ),
                ),
                const SizedBox(width: 12),
                SizedBox(
                  height: 52,
                  child: ElevatedButton(
                    onPressed: _isLoading ? null : _search,
                    child: _isLoading
                        ? const SizedBox(
                            width: 20,
                            height: 20,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              color: Colors.white,
                            ),
                          )
                        : const Icon(Icons.search),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 24),

            // Error
            if (_errorMessage != null)
              Card(
                color: AppTheme.errorColor.withValues(alpha: 0.1),
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Row(
                    children: [
                      const Icon(
                        Icons.error_outline,
                        color: AppTheme.errorColor,
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Text(
                          _errorMessage!,
                          style: const TextStyle(color: AppTheme.errorColor),
                        ),
                      ),
                    ],
                  ),
                ),
              ),

            // Client info
            if (_client != null) ...[
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(20),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Container(
                            padding: const EdgeInsets.all(10),
                            decoration: BoxDecoration(
                              color: AppTheme.primaryColor.withValues(
                                alpha: 0.1,
                              ),
                              borderRadius: BorderRadius.circular(10),
                            ),
                            child: const Icon(
                              Icons.person,
                              color: AppTheme.primaryColor,
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  '${_client!['nombres']} ${_client!['apellidos']}',
                                  style: const TextStyle(
                                    fontSize: 18,
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                                Text(
                                  'Cédula: ${CedulaUtils.format(_client!['cedula'])}',
                                  style: TextStyle(
                                    color: AppTheme.textSecondary,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                      const Divider(height: 24),
                      // Balance
                      Container(
                        padding: const EdgeInsets.all(16),
                        decoration: BoxDecoration(
                          color: AppTheme.primaryColor.withValues(alpha: 0.08),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Icon(
                              Icons.stars_outlined,
                              color: AppTheme.primaryColor,
                            ),
                            const SizedBox(width: 8),
                            Text(
                              '${_balance?.toStringAsFixed(0) ?? '0'} puntos',
                              style: const TextStyle(
                                fontSize: 22,
                                fontWeight: FontWeight.w700,
                                color: AppTheme.primaryColor,
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 16),
                      // Actions
                      Row(
                        children: [
                          Expanded(
                            child: ElevatedButton.icon(
                              onPressed: () {
                                context.push(
                                  '/transactions/earn',
                                  extra: {
                                    'clientId': _client!['id'],
                                    'clientName':
                                        '${_client!['nombres']} ${_client!['apellidos']}',
                                    'clientCedula': _client!['cedula'],
                                  },
                                );
                              },
                              icon: const Icon(Icons.add_circle, size: 20),
                              label: const Text('Acumular'),
                              style: ElevatedButton.styleFrom(
                                backgroundColor: AppTheme.successColor,
                              ),
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: ElevatedButton.icon(
                              onPressed: () {
                                context.push(
                                  '/transactions/redeem',
                                  extra: {
                                    'clientId': _client!['id'],
                                    'clientName':
                                        '${_client!['nombres']} ${_client!['apellidos']}',
                                    'clientCedula': _client!['cedula'],
                                    'balance': _balance ?? 0.0,
                                    'minRedeem': _minRedeem,
                                  },
                                );
                              },
                              icon: const Icon(Icons.redeem_outlined, size: 20),
                              label: const Text('Canjear'),
                              style: ElevatedButton.styleFrom(
                                backgroundColor: AppTheme.primaryDark,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),

              // Historial de transacciones
              const SizedBox(height: 24),
              Text(
                'Historial',
                style: const TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w600,
                  color: AppTheme.textSecondary,
                  letterSpacing: 0.8,
                ),
              ),
              const SizedBox(height: 12),

              if (_history.isEmpty)
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(24),
                    child: Center(
                      child: Text(
                        'Sin transacciones',
                        style: TextStyle(
                          color: AppTheme.textSecondary,
                          fontSize: 14,
                        ),
                      ),
                    ),
                  ),
                )
              else
                ..._history.map((tx) {
                  final type = tx['transaction_type'] as String;
                  final isEarn = type == 'earn';
                  final points = (tx['points'] as num).toDouble();
                  final gallons =
                      (tx['gallons_amount'] as num?)?.toDouble() ?? 0;
                  final createdAt = tx['created_at'] as String;
                  final expiresAt = tx['expires_at'] as String?;
                  final expired = _isExpired(expiresAt);

                  return Padding(
                    padding: const EdgeInsets.only(bottom: 8),
                    child: Card(
                      child: Padding(
                        padding: const EdgeInsets.all(14),
                        child: Row(
                          children: [
                            // Icono tipo
                            Container(
                              width: 36,
                              height: 36,
                              decoration: BoxDecoration(
                                color: isEarn
                                    ? AppTheme.successColor.withValues(
                                        alpha: 0.1,
                                      )
                                    : AppTheme.primaryColor.withValues(
                                        alpha: 0.1,
                                      ),
                                borderRadius: BorderRadius.circular(8),
                              ),
                              child: Icon(
                                isEarn
                                    ? Icons.add_circle_outline
                                    : Icons.remove_circle_outline,
                                size: 20,
                                color: isEarn
                                    ? AppTheme.successColor
                                    : AppTheme.primaryColor,
                              ),
                            ),
                            const SizedBox(width: 12),
                            // Info
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    isEarn ? 'Acumulación' : 'Canje',
                                    style: const TextStyle(
                                      fontSize: 14,
                                      fontWeight: FontWeight.w600,
                                      color: AppTheme.textPrimary,
                                    ),
                                  ),
                                  const SizedBox(height: 2),
                                  Text(
                                    _formatDate(createdAt),
                                    style: const TextStyle(
                                      fontSize: 12,
                                      color: AppTheme.textSecondary,
                                    ),
                                  ),
                                  if (isEarn && gallons > 0)
                                    Text(
                                      '${gallons.toStringAsFixed(1)} galones',
                                      style: const TextStyle(
                                        fontSize: 12,
                                        color: AppTheme.textSecondary,
                                      ),
                                    ),
                                  if (expiresAt != null)
                                    Text(
                                      'Expira: ${_formatExpiryDate(expiresAt)}',
                                      style: TextStyle(
                                        fontSize: 11,
                                        color: expired
                                            ? AppTheme.errorColor
                                            : AppTheme.textSecondary,
                                        fontWeight: expired
                                            ? FontWeight.w600
                                            : FontWeight.normal,
                                      ),
                                    ),
                                ],
                              ),
                            ),
                            // Points + Reprint
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.end,
                              children: [
                                Text(
                                  '${isEarn ? '+' : ''}${points.toStringAsFixed(0)}',
                                  style: TextStyle(
                                    fontSize: 16,
                                    fontWeight: FontWeight.w700,
                                    color: isEarn
                                        ? AppTheme.successColor
                                        : AppTheme.primaryColor,
                                  ),
                                ),
                                const SizedBox(height: 4),
                                GestureDetector(
                                  onTap: () => _reprintTransaction(tx),
                                  child: Icon(
                                    Icons.print_outlined,
                                    size: 18,
                                    color: AppTheme.textSecondary,
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                  );
                }),
            ],
          ],
        ),
      ),
    );
  }
}
