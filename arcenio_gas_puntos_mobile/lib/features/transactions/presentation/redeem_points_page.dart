import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../../../core/network/dio_client.dart';
import '../../../core/printer/printer_service.dart';
import '../../../core/printer/ticket_generator.dart';
import '../../../core/theme.dart';
import '../data/transaction_remote_datasource.dart';

class RedeemPointsPage extends StatefulWidget {
  final String clientId;
  final String clientName;
  final String clientCedula;
  final double currentBalance;
  final double minRedeem;

  const RedeemPointsPage({
    super.key,
    required this.clientId,
    required this.clientName,
    required this.clientCedula,
    required this.currentBalance,
    required this.minRedeem,
  });

  @override
  State<RedeemPointsPage> createState() => _RedeemPointsPageState();
}

class _RedeemPointsPageState extends State<RedeemPointsPage> {
  bool _isLoading = false;
  String? _errorMessage;
  Map<String, dynamic>? _result;

  Future<void> _redeemAll() async {
    setState(() {
      _isLoading = true;
      _errorMessage = null;
      _result = null;
    });

    try {
      final dioClient = context.read<DioClient>();
      final datasource = TransactionRemoteDatasource(dioClient);
      final result = await datasource.redeemPoints(widget.clientId);

      if (mounted) {
        setState(() => _result = result);
        // Auto-imprimir
        _autoPrint(result);
      }
    } catch (e) {
      setState(() {
        _errorMessage = e.toString().replaceAll('Exception: ', '');
      });
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Future<void> _autoPrint(Map<String, dynamic> result) async {
    try {
      final bytes = await TicketGenerator.generateRedeemTicket(
        clientName: widget.clientName,
        clientCedula: widget.clientCedula,
        pointsRedeemed: (result['points_redeemed'] as num).toDouble(),
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
    await _autoPrint(_result!);
    if (mounted) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('Reimprimiendo ticket...')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Canjear Puntos')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Client info + balance
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  children: [
                    Row(
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
                    const Divider(height: 20),
                    Container(
                      padding: const EdgeInsets.all(12),
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
                            '${widget.currentBalance.toStringAsFixed(0)} puntos disponibles',
                            style: const TextStyle(
                              fontSize: 17,
                              fontWeight: FontWeight.w700,
                              color: AppTheme.primaryColor,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 20),

            if (_result == null) ...[
              // Mensaje de confirmación
              Text(
                'Se canjearán todos los puntos disponibles del cliente.\n'
                'Mínimo requerido: ${widget.minRedeem.toStringAsFixed(0)} puntos.',
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 14,
                  color: AppTheme.textSecondary,
                  height: 1.5,
                ),
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
                  onPressed: _isLoading ? null : _redeemAll,
                  icon: _isLoading
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            color: Colors.white,
                          ),
                        )
                      : const Icon(Icons.redeem_outlined, size: 20),
                  label: Text(
                    'Canjear ${widget.currentBalance.toStringAsFixed(0)} Puntos',
                  ),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppTheme.primaryDark,
                  ),
                ),
              ),
            ],

            if (_result != null) ...[
              Card(
                color: AppTheme.primaryDark.withValues(alpha: 0.04),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                  side: BorderSide(
                    color: AppTheme.primaryDark.withValues(alpha: 0.25),
                  ),
                ),
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    children: [
                      const Icon(
                        Icons.check_circle_outline,
                        color: AppTheme.primaryDark,
                        size: 48,
                      ),
                      const SizedBox(height: 8),
                      const Text(
                        'Puntos Canjeados',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w700,
                          color: AppTheme.primaryDark,
                        ),
                      ),
                      const SizedBox(height: 20),
                      _ResultRow(
                        label: 'Puntos canjeados',
                        value:
                            '-${(_result!['points_redeemed'] as num).toStringAsFixed(0)}',
                        valueColor: AppTheme.errorColor,
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
