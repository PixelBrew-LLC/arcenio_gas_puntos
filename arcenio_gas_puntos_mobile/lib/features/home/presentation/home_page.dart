import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:toastification/toastification.dart';
import '../../../core/printer/printer_service.dart';
import '../../../core/storage/secure_storage_service.dart';
import '../../../core/theme.dart';
import 'package:provider/provider.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  String _userName = '';
  String _userRole = '';

  @override
  void initState() {
    super.initState();
    _loadUserData();
    _tryConnectPrinter();
  }

  Future<void> _tryConnectPrinter() async {
    final printerService = context.read<PrinterService>();

    // If no printer saved, nothing to connect to
    if (printerService.savedPrinter == null) return;

    // If already connected, show success immediately
    if (printerService.isConnected) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (mounted) _showPrinterToast(success: true);
      });
      return;
    }

    // Attempt connection and wait a bit for the BT status to settle
    await printerService.connectToSaved();
    await Future.delayed(const Duration(seconds: 3));

    if (mounted) {
      _showPrinterToast(success: printerService.isConnected);
    }
  }

  void _showPrinterToast({required bool success}) {
    toastification.show(
      context: context,
      type: success ? ToastificationType.success : ToastificationType.error,
      style: ToastificationStyle.flatColored,
      title: Text(success ? 'Impresora conectada' : 'Sin conexión a impresora'),
      description: Text(
        success
            ? 'La impresora se conectó exitosamente.'
            : 'No se pudo conectar a la impresora.',
      ),
      alignment: Alignment.topCenter,
      autoCloseDuration: const Duration(seconds: 4),
      showProgressBar: false,
    );
  }

  Future<void> _loadUserData() async {
    final storageService = context.read<SecureStorageService>();
    final userData = await storageService.getUserData();
    if (userData != null && mounted) {
      setState(() {
        _userName =
            '${userData['nombres'] ?? ''} ${userData['apellidos'] ?? ''}'
                .trim();
        _userRole = userData['role'] ?? '';
      });
    }
  }

  Future<void> _logout() async {
    final storageService = context.read<SecureStorageService>();
    await storageService.clearAll();
    if (mounted) {
      context.go('/login');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Arcenio Gas'),
        actions: [
          IconButton(
            icon: const Icon(Icons.print_outlined),
            onPressed: () => context.push('/settings/printer'),
            tooltip: 'Impresora',
          ),
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: _logout,
            tooltip: 'Cerrar Sesión',
          ),
        ],
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Bienvenida
              Row(
                children: [
                  Container(
                    width: 44,
                    height: 44,
                    decoration: BoxDecoration(
                      color: AppTheme.primaryColor,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: const Icon(
                      Icons.person,
                      color: Colors.white,
                      size: 22,
                    ),
                  ),
                  const SizedBox(width: 14),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          _userName.isEmpty ? 'Bienvenido' : _userName,
                          style: const TextStyle(
                            fontSize: 17,
                            fontWeight: FontWeight.w700,
                            color: AppTheme.textPrimary,
                          ),
                        ),
                        if (_userRole.isNotEmpty)
                          Text(
                            _userRole,
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
              const SizedBox(height: 28),

              const Text(
                'Acciones',
                style: TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w600,
                  color: AppTheme.textSecondary,
                  letterSpacing: 0.8,
                ),
              ),
              const SizedBox(height: 16),

              // Acciones
              IntrinsicHeight(
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Expanded(
                      child: _ActionCard(
                        icon: Icons.search,
                        label: 'Buscar Cliente',
                        color: AppTheme.primaryColor,
                        onTap: () => context.push('/clients/search'),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: _ActionCard(
                        icon: Icons.person_add_outlined,
                        label: 'Registrar Cliente',
                        color: AppTheme.successColor,
                        onTap: () => context.push('/clients/register'),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _ActionCard extends StatelessWidget {
  final IconData icon;
  final String label;
  final Color color;
  final VoidCallback onTap;

  const _ActionCard({
    required this.icon,
    required this.label,
    required this.color,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Colors.white,
      borderRadius: BorderRadius.circular(8),
      clipBehavior: Clip.hardEdge,
      child: InkWell(
        onTap: onTap,
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            border: Border.all(color: AppTheme.dividerColor),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: color.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Icon(icon, size: 22, color: color),
              ),
              const SizedBox(height: 12),
              Text(
                label,
                style: const TextStyle(
                  color: AppTheme.textPrimary,
                  fontSize: 13,
                  fontWeight: FontWeight.w600,
                  height: 1.3,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
