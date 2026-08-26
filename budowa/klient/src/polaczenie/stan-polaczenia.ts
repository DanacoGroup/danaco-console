/**
 * Stan warstwy połączenia widziany przez pozostałe warstwy klienta.
 *
 * Stan jest wyłącznie informacją: rozłączenie nie odbiera warstwom wyższym
 * prawa do wysyłania. Ramka wpisana przy rozłączeniu trafia do kolejki
 * wychodzącej i idzie do rdzenia po wznowieniu połączenia.
 */
export type StanPolaczenia = 'rozlaczony' | 'laczenie' | 'polaczony' | 'ponawianie';
