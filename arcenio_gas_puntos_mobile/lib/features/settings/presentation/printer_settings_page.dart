import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../core/printer/printer_service.dart';
import '../../../core/theme.dart';

class PrinterSettingsPage extends StatefulWidget {
  const PrinterSettingsPage({super.key});

  @override
  State<PrinterSettingsPage> createState() => _PrinterSettingsPageState();
}

class _PrinterSettingsPageState extends State<PrinterSettingsPage> {
  @override
  void initState() {
    super.initState();
    // Iniciar escaneo al abrir
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<PrinterService>().scan();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Impresora')),
      body: Consumer<PrinterService>(
        builder: (context, printer, _) {
          return Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Estado actual
                _ConnectionStatus(
                  isConnected: printer.isConnected,
                  printerName: printer.savedPrinter?.name,
                ),
                const SizedBox(height: 16),

                // Botones de acción
                if (printer.savedPrinter != null)
                  SizedBox(
                    width: double.infinity,
                    child: OutlinedButton.icon(
                      onPressed: () => printer.forgetPrinter(),
                      icon: const Icon(Icons.link_off, size: 20),
                      label: const Text('Olvidar Impresora'),
                      style: OutlinedButton.styleFrom(
                        foregroundColor: AppTheme.errorColor,
                        side: const BorderSide(color: AppTheme.errorColor),
                        padding: const EdgeInsets.symmetric(vertical: 14),
                      ),
                    ),
                  ),

                const SizedBox(height: 24),

                // Sección de escaneo
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text(
                      'Dispositivos',
                      style: TextStyle(
                        fontSize: 13,
                        fontWeight: FontWeight.w600,
                        color: AppTheme.textSecondary,
                        letterSpacing: 0.8,
                      ),
                    ),
                    TextButton.icon(
                      onPressed: printer.isScanning
                          ? null
                          : () => printer.scan(),
                      icon: printer.isScanning
                          ? const SizedBox(
                              width: 16,
                              height: 16,
                              child: CircularProgressIndicator(strokeWidth: 2),
                            )
                          : const Icon(Icons.refresh, size: 18),
                      label: Text(
                        printer.isScanning ? 'Buscando...' : 'Buscar',
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 8),

                // Lista de dispositivos
                Expanded(
                  child: printer.discoveredDevices.isEmpty
                      ? Center(
                          child: Text(
                            printer.isScanning
                                ? 'Buscando impresoras...'
                                : 'No se encontraron impresoras.\nVerifique que el Bluetooth esté encendido.',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              color: AppTheme.textSecondary,
                              fontSize: 14,
                            ),
                          ),
                        )
                      : ListView.separated(
                          itemCount: printer.discoveredDevices.length,
                          separatorBuilder: (context, index) =>
                              const SizedBox(height: 8),
                          itemBuilder: (context, index) {
                            final device = printer.discoveredDevices[index];
                            final isSelected =
                                printer.savedPrinter?.address == device.address;

                            return Card(
                              color: isSelected
                                  ? AppTheme.primaryColor.withValues(
                                      alpha: 0.05,
                                    )
                                  : null,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(8),
                                side: BorderSide(
                                  color: isSelected
                                      ? AppTheme.primaryColor
                                      : AppTheme.dividerColor,
                                ),
                              ),
                              child: ListTile(
                                contentPadding: const EdgeInsets.symmetric(
                                  horizontal: 16,
                                  vertical: 4,
                                ),
                                leading: Icon(
                                  Icons.print_outlined,
                                  color: isSelected
                                      ? AppTheme.primaryColor
                                      : AppTheme.textSecondary,
                                ),
                                title: Text(
                                  device.name.isNotEmpty
                                      ? device.name
                                      : 'Impresora desconocida',
                                  style: TextStyle(
                                    fontWeight: isSelected
                                        ? FontWeight.w600
                                        : FontWeight.normal,
                                    fontSize: 14,
                                  ),
                                ),
                                subtitle: Text(
                                  device.address ?? '',
                                  style: const TextStyle(fontSize: 12),
                                ),
                                trailing: isSelected
                                    ? Icon(
                                        printer.isConnected
                                            ? Icons.bluetooth_connected
                                            : Icons.check_circle_outline,
                                        color: printer.isConnected
                                            ? AppTheme.successColor
                                            : AppTheme.primaryColor,
                                      )
                                    : null,
                                onTap: () => printer.selectAndConnect(device),
                              ),
                            );
                          },
                        ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}

class _ConnectionStatus extends StatelessWidget {
  final bool isConnected;
  final String? printerName;

  const _ConnectionStatus({required this.isConnected, this.printerName});

  @override
  Widget build(BuildContext context) {
    final hasPrinter = printerName != null;

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isConnected
            ? AppTheme.successColor.withValues(alpha: 0.08)
            : AppTheme.surfaceVariant,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(
          color: isConnected
              ? AppTheme.successColor.withValues(alpha: 0.3)
              : AppTheme.dividerColor,
        ),
      ),
      child: Row(
        children: [
          Icon(
            isConnected
                ? Icons.bluetooth_connected
                : hasPrinter
                ? Icons.bluetooth_searching
                : Icons.bluetooth_disabled,
            color: isConnected ? AppTheme.successColor : AppTheme.textSecondary,
            size: 28,
          ),
          const SizedBox(width: 14),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  isConnected
                      ? 'Conectada'
                      : hasPrinter
                      ? 'Desconectada'
                      : 'Sin impresora',
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w600,
                    color: isConnected
                        ? AppTheme.successColor
                        : AppTheme.textPrimary,
                  ),
                ),
                if (hasPrinter)
                  Text(
                    printerName!,
                    style: const TextStyle(
                      fontSize: 13,
                      color: AppTheme.textSecondary,
                    ),
                  ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
