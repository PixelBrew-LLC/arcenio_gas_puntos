import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:toastification/toastification.dart';
import 'core/network/dio_client.dart';
import 'core/printer/printer_service.dart';
import 'core/router/app_router.dart';
import 'core/storage/secure_storage_service.dart';
import 'core/theme.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(const ArcenioGasApp());
}

class ArcenioGasApp extends StatelessWidget {
  const ArcenioGasApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        Provider<SecureStorageService>(create: (_) => SecureStorageService()),
        ChangeNotifierProvider<PrinterService>(
          create: (_) {
            final service = PrinterService();
            service.init();
            return service;
          },
        ),
        ProxyProvider<SecureStorageService, DioClient>(
          update: (_, storageService, prev) => DioClient(storageService),
        ),
        ProxyProvider<SecureStorageService, AppRouter>(
          update: (_, storageService, prev) => AppRouter(storageService),
        ),
      ],
      child: Builder(
        builder: (context) {
          final appRouter = context.read<AppRouter>();
          return ToastificationWrapper(
            child: MaterialApp.router(
              title: 'Arcenio Gas Puntos',
              debugShowCheckedModeBanner: false,
              theme: AppTheme.theme,
              routerConfig: appRouter.router,
            ),
          );
        },
      ),
    );
  }
}
