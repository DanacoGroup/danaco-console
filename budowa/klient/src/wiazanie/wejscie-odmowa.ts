// Odmowa drogi wejścia czytana z kodu kontraktu. Rdzeń dokłada do `details`
// powód maszynowy, bo sam kod bywa za szeroki — `conflict` niesie i login
// zajęty, i adres potwierdzony wcześniej.
import { ErrorCode, type ErrorInfo } from '../../../shared/contract.ts';
import { oglos } from './ogloszenie.ts';
import { napis } from './wejscie-katalog.ts';

/** Powody maszynowe rdzenia z `core/adapter_modul_auth_stan.go`. */
export const PowodOdmowy = {
  KontoIstnieje: 'konto-istnieje',
  KolizjaDanych: 'kolizja-danych',
  AdresNiepotwierdzony: 'adres-niepotwierdzony',
} as const;

/** Zdania zapasowe na wypadek odmowy bez opisu; kolejność jak w katalogu kodów kontraktu. */
const ZDANIA: Readonly<Record<ErrorCode, string>> = {
  [ErrorCode.ValidationFailed]: 'Treść żądania nie zgadza się z kontraktem.',
  [ErrorCode.NotFound]: 'Rdzeń nie zna wskazanego bytu.',
  [ErrorCode.NotAuthenticated]: 'Rdzeń nie rozpoznał danych wejścia.',
  [ErrorCode.PermissionDenied]: 'Konto nie ma uprawnienia do tej czynności.',
  [ErrorCode.Conflict]: 'Stan konta wyklucza tę czynność.',
  [ErrorCode.ChannelUnavailable]: 'Kanał modelu nie odpowiada.',
  [ErrorCode.RateLimited]: 'Rdzeń ogranicza tempo żądań.',
  [ErrorCode.InternalError]: 'Rdzeń zgłosił błąd wewnętrzny.',
};

/** Kod odmowy z kontraktu; pustka znaczy niepowodzenie bez odpowiedzi rdzenia. */
export function kodOdmowy(blad: ErrorInfo | undefined): ErrorCode | '' {
  return blad?.code ?? '';
}

/** Powód maszynowy dołożony przez rdzeń; pustka, gdy odmowa go nie niesie. */
export function powodOdmowy(blad: ErrorInfo | undefined): string {
  const dane = blad?.details as { powod?: unknown } | undefined;
  return typeof dane?.powod === 'string' ? dane.powod : '';
}

/** Zdanie odmowy: opis rdzenia, a gdy go brak — zdanie zapasowe kodu kontraktu. */
export function zdanieOdmowy(blad: ErrorInfo | undefined): string {
  const opis = blad?.message.trim() ?? '';
  if (opis !== '') return opis;
  const kod = kodOdmowy(blad);
  return kod === '' ? 'Rdzeń nie odpowiedział na żądanie.' : ZDANIA[kod];
}

/**
 * Stawia odmowę przed Operatorem: w banerze odsłony, a gdy odsłona banera nie
 * ma — komunikatem. Baner niesie znacznik `data-usterka-formularza`, więc
 * sprzątnie go to samo sprawdzenie pól, które sprząta usterki wpisu.
 */
export function pokazOdmowe(
  panel: HTMLElement | null,
  glowa: string,
  blad: ErrorInfo | undefined,
): void {
  const tresc = zdanieOdmowy(blad);
  if (panel === null || !wpiszWBaner(panel, glowa, tresc)) oglos(glowa, tresc, 'blad');
}

/** Odsłona logowania po odmowie; adres niepotwierdzony to konto zatrzymane przed aktywacją, nie zła para login–hasło. */
export function odmowaLogowania(blad: ErrorInfo | undefined): void {
  const przelacz = (globalThis as {
    dnPrzelaczWidok?: (widok: string, grupa?: string | null) => void;
  }).dnPrzelaczWidok;
  if (powodOdmowy(blad) === PowodOdmowy.AdresNiepotwierdzony) {
    przelacz?.('kod', 'stan');
    return;
  }
  przelacz?.('logowanie-blad', 'stan');
  const panel = document.querySelector<HTMLElement>('.we-panel[data-widok="logowanie-blad"]');
  pokazOdmowe(panel, napis('usterki.naglowekLogowanie'), blad);
}

/**
 * Odmowa rejestracji. Login albo adres zajęty wraca kodem `conflict` z powodem
 * `kolizja-danych`, więc odsłona nazywa rzecz zdaniem katalogu; pozostałe kody
 * idą zdaniem rdzenia.
 */
export function odmowaRejestracji(panel: HTMLElement, blad: ErrorInfo | undefined): void {
  if (kodOdmowy(blad) === ErrorCode.Conflict && powodOdmowy(blad) === PowodOdmowy.KolizjaDanych) {
    pokazOdmowe(panel, napis('usterki.loginZajety.glowa'), {
      code: ErrorCode.Conflict,
      message: napis('usterki.loginZajety.tresc'),
      retryable: false,
    });
    return;
  }
  pokazOdmowe(panel, napis('usterki.naglowekKonto'), blad);
}

/** Wpisuje odmowę w baner odsłony; fałsz znaczy odsłonę bez banera. */
function wpiszWBaner(panel: HTMLElement, glowa: string, tresc: string): boolean {
  const baner = panel.querySelector<HTMLElement>('.we-komunikaty .dn-alert');
  const pole = baner?.querySelector<HTMLElement>('.dn-alert-tresc');
  const czolo = pole?.querySelector('b');
  if (baner == null || pole == null || czolo == null) return false;
  baner.classList.remove('dn-alert--info', 'dn-alert--ostrzezenie', 'dn-alert--sukces');
  baner.classList.add('dn-alert--blad');
  baner.setAttribute('role', 'alert');
  czolo.textContent = glowa;
  pole.querySelector('.dn-alert-lista')?.remove();
  const ostatni = pole.lastChild;
  if (ostatni !== null && ostatni.nodeType === Node.TEXT_NODE) ostatni.textContent = tresc;
  else pole.appendChild(document.createTextNode(tresc));
  return true;
}
