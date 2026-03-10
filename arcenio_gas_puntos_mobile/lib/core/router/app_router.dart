import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../features/auth/presentation/login_page.dart';
import '../../features/home/presentation/home_page.dart';
import '../../features/clients/presentation/client_search_page.dart';
import '../../features/clients/presentation/client_register_page.dart';
import '../../features/transactions/presentation/earn_points_page.dart';
import '../../features/transactions/presentation/redeem_points_page.dart';
import '../../features/settings/presentation/printer_settings_page.dart';
import '../storage/secure_storage_service.dart';

class AppRouter {
  final SecureStorageService storageService;

  AppRouter(this.storageService);

  late final GoRouter router = GoRouter(
    initialLocation: '/login',
    redirect: (BuildContext context, GoRouterState state) async {
      final isLoggedIn = await storageService.isTokenValidForToday();
      final isLoginRoute = state.matchedLocation == '/login';

      if (!isLoggedIn && !isLoginRoute) {
        return '/login';
      }

      if (isLoggedIn && isLoginRoute) {
        return '/home';
      }

      return null;
    },
    routes: [
      GoRoute(
        path: '/login',
        name: 'login',
        builder: (context, state) => const LoginPage(),
      ),
      GoRoute(
        path: '/home',
        name: 'home',
        builder: (context, state) => const HomePage(),
      ),
      GoRoute(
        path: '/clients/search',
        name: 'clientSearch',
        builder: (context, state) => const ClientSearchPage(),
      ),
      GoRoute(
        path: '/clients/register',
        name: 'clientRegister',
        builder: (context, state) {
          final extra = state.extra as Map<String, dynamic>?;
          return ClientRegisterPage(
            initialCedula: extra?['initialCedula'] as String?,
          );
        },
      ),
      GoRoute(
        path: '/transactions/earn',
        name: 'earnPoints',
        builder: (context, state) {
          final extra = state.extra as Map<String, dynamic>?;
          return EarnPointsPage(
            clientId: extra?['clientId'] ?? '',
            clientName: extra?['clientName'] ?? '',
            clientCedula: extra?['clientCedula'] ?? '',
          );
        },
      ),
      GoRoute(
        path: '/transactions/redeem',
        name: 'redeemPoints',
        builder: (context, state) {
          final extra = state.extra as Map<String, dynamic>?;
          return RedeemPointsPage(
            clientId: extra?['clientId'] ?? '',
            clientName: extra?['clientName'] ?? '',
            clientCedula: extra?['clientCedula'] ?? '',
            currentBalance: extra?['balance'] ?? 0.0,
            minRedeem: extra?['minRedeem'] ?? 0.0,
          );
        },
      ),
      GoRoute(
        path: '/settings/printer',
        name: 'printerSettings',
        builder: (context, state) => const PrinterSettingsPage(),
      ),
    ],
  );
}
