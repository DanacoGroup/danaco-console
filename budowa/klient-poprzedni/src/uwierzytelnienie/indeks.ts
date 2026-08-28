/** Moduł uwierzytelnienia wystawia bramkę operatora publicznie jako ekran logowania punktu wejścia oraz trwałość sesji bramki dla warstw, które po wejściu potrzebują tokenu. */
export { utworzEkranLogowania } from './ekran-logowania';
export type { EkranLogowania, OpisEkranuLogowania } from './ekran-logowania';
export { odczytajSesje, opisWaznosci, sesjaTrwala } from './sesja-bramki';
/** Tożsamość maszyny — sam odczyt, bo bramka identyfikatora nie nadaje: powstaje on przy zakładaniu PIN-u w oknie ustawień. */
export { odczytajTozsamoscUrzadzenia, KLUCZ_URZADZENIA } from './tozsamosc-urzadzenia';
