/**
 * Stan warstwy połączenia widziany przez pozostałe warstwy klienta. Stan jest
 * wyłącznie informacją: rozłączenie nie odbiera warstwom wyższym prawa do
 * wysyłania.
 */
export type StanPolaczenia = 'rozlaczony' | 'laczenie' | 'polaczony' | 'ponawianie';
