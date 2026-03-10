import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../../../core/network/dio_client.dart';
import '../../../core/printer/printer_service.dart';
import '../../../core/printer/ticket_generator.dart';
import '../../../core/theme.dart';
import '../data/transaction_remote_datasource.dart';

class EarnPointsPage extends StatefulWidget {
  final String clientId;
  final String clientName;
  final String clientCedula;

  const EarnPointsPage({
    super.key,
    required this.clientId,
    required this.clientName,
    required this.clientCedula,
  });

  @override
  State<EarnPointsPage> createState() => _EarnPointsPageState();
}

class _EarnPointsPageState extends State<EarnPointsPage> {
  final _gallonsController = TextEditingController();
  bool _isLoading = false;
  String? _errorMessage;
  Map<String, dynamic>? _result;

  @override
  void dispose() {
    _gallonsController.dispose();
    super.dispose();
  }

  Future<void> _earn() async {
    final gallons = double.tryParse(_gallonsController.text);
    if (gallons == null || gallons <= 0) {
      setState(() => _errorMessage = 'Ingrese una cantidad válida de galones');
      return;
    }

    setState(() {
      _isLoading = true;
      _errorMessage = null;
      _result = null;
    });

    try {
      final dioClient = context.read<DioClient>();
      final datasource = TransactionRemoteDatasource(dioClient);
      final result = await datasource.earnPoints(widget.clientId, gallons);

      if (mounted) {
        setState(() => _result = result);
        // Auto-imprimir
        _autoPrint(result, gallons);
      }
    } catch (e) {
      setState(() {
        _errorMessage = e.toString().replaceAll('Exception: ', '');
      });
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Future<void> _autoPrint(Map<String, dynamic> result, double gallons) async {
    try {
      final bytes = await TicketGenerator.generateEarnTicket(
        clientName: widget.clientName,
        clientCedula: widget.clientCedula,
        gallons: gallons,
        pointsEarned: (result['points_earned'] as num).toDouble(),
        newBalance: (result['new_balance'] as num).toDouble(),
        date: DateTime.now(),
      );
      if (mounted) {
        await context.read<PrinterService>().printBytes(bytes);
      }
    } catch (_) {
      // Si falla la impresión no bloquear el flujo
    }
  }

  Future<void> _reprint() async {
    if (_result == null) return;
    final gallons = double.tryParse(_gallonsController.text) ?? 0;
    await _autoPrint(_result!, gallons);
    if (mounted) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('Reimprimiendo ticket...')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Acumular Puntos')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Client info
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Row(
                  children: [
                    const Icon(Icons.person, color: AppTheme.primaryColor),
                    const SizedBox(width: 12),
                    Text(
                      widget.clientName,
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 20),

            if (_result == null) ...[
              // Input
              TextField(
                controller: _gallonsController,
                decoration: const InputDecoration(
                  labelText: 'Cantidad de Galones',
                  prefixIcon: Icon(Icons.local_gas_station),
                  hintText: 'Ej: 15.5',
                ),
                keyboardType: const TextInputType.numberWithOptions(
                  decimal: true,
                ),
                autofocus: true,
              ),
              const SizedBox(height: 24),

              if (_errorMessage != null)
                Padding(
                  padding: const EdgeInsets.only(bottom: 16),
                  child: Text(
                    _errorMessage!,
                    textAlign: TextAlign.center,
                    style: const TextStyle(color: AppTheme.errorColor),
                  ),
                ),

              SizedBox(
                height: 50,
                child: ElevatedButton.icon(
                  onPressed: _isLoading ? null : _earn,
                  icon: _isLoading
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            color: Colors.white,
                          ),
                        )
                      : const Icon(Icons.add_circle_outline, size: 20),
                  label: const Text('Acumular Puntos'),
                ),
              ),
            ],

            // Result
            if (_result != null) ...[
              Card(
                color: AppTheme.successColor.withValues(alpha: 0.04),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                  side: BorderSide(
                    color: AppTheme.successColor.withValues(alpha: 0.25),
                  ),
                ),
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    children: [
                      const Icon(
                        Icons.check_circle_outline,
                        color: AppTheme.successColor,
                        size: 48,
                      ),
                      const SizedBox(height: 8),
                      const Text(
                        'Puntos Acumulados',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w700,
                          color: AppTheme.successColor,
                        ),
                      ),
                      const SizedBox(height: 20),
                      _ResultRow(
                        label: 'Puntos ganados',
                        value:
                            '+${(_result!['points_earned'] as num).toStringAsFixed(0)}',
                        valueColor: AppTheme.successColor,
                      ),
                      const SizedBox(height: 8),
                      _ResultRow(
                        label: 'Nuevo saldo',
                        value:
                            '${(_result!['new_balance'] as num).toStringAsFixed(0)} pts',
                        valueColor: AppTheme.primaryColor,
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 20),
              OutlinedButton.icon(
                onPressed: _reprint,
                icon: const Icon(Icons.print),
                label: const Text('Reimprimir Ticket'),
                style: OutlinedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 14),
                ),
              ),
              const SizedBox(height: 12),
              ElevatedButton(
                onPressed: () => context.go('/home'),
                child: const Text('Volver al Inicio'),
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _ResultRow extends StatelessWidget {
  final String label;
  final String value;
  final Color valueColor;

  const _ResultRow({
    required this.label,
    required this.value,
    required this.valueColor,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: TextStyle(color: AppTheme.textSecondary, fontSize: 15),
        ),
        Text(
          value,
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.w700,
            color: valueColor,
          ),
        ),
      ],
    );
  }
}
