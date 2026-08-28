/**
 * Stan warstwy transportu widziany przez pozostałe warstwy klienta jest wyłącznie informacją: okno pozostaje użyteczne w każdym stanie, a wiadomości wpisane przy rozłączeniu trafiają do kolejki wychodzącej.
 */
export type StanPolaczenia = 'rozlaczony' | 'laczenie' | 'polaczony' | 'ponawianie';
