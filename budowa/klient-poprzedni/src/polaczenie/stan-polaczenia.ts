/**
 * Stan warstwy transportu widziany przez pozostałe warstwy klienta.
 *
 * Stan jest wyłącznie informacją i nie służy do wyłączania elementów
 * interfejsu — okno komunikacji pozostaje użyteczne w każdym stanie,
 * a wiadomości wpisane przy rozłączeniu trafiają do kolejki wychodzącej.
 */
export type StanPolaczenia = 'rozlaczony' | 'laczenie' | 'polaczony' | 'ponawianie';
