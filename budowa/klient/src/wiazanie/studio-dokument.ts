// Dokument Studia: odczyt dokumentu okna, zapis treści wraz z postacią
// i format dokumentu wskazany w pasie stanu kanwy.

import {
  Command,
  ErrorCode,
  EventType,
  StudioDocumentFormat,
  type ErrorInfo,
  type StudioActionBalance,
  type StudioDocument,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const POWODY_W_KOMUNIKACIE = 3;

const ZDANIA_ODMOWY: Readonly<Record<ErrorCode, string>> = {
  [ErrorCode.ValidationFailed]: 'żądanie nie zgadza się z kontraktem',
  [ErrorCode.NotFound]: 'rdzeń nie zna wskazanego wiersza',
  [ErrorCode.NotAuthenticated]: 'sesja wejścia wygasła',
  [ErrorCode.PermissionDenied]: 'konto nie ma uprawnienia do tej czynności',
  [ErrorCode.Conflict]: 'inna zmiana wyprzedziła tę czynność',
  [ErrorCode.ChannelUnavailable]: 'kanał modelu nie odpowiada',
  [ErrorCode.RateLimited]: 'rdzeń ogranicza tempo żądań',
  [ErrorCode.InternalError]: 'rdzeń zgłosił błąd wewnętrzny',
};

export function nazwijOdmowe(czynnosc: string, blad: ErrorInfo | undefined): string {
  const kod = blad?.code;
  const powod = kod === undefined ? 'rdzeń nie odpowiedział' : ZDANIA_ODMOWY[kod];
  return `${czynnosc} — ${powod}.`;
}

export function nazwijBilans(czynnosc: string, bilans: StudioActionBalance): string {
  if (bilans.skippedCount === 0) return `${czynnosc} — zmieniono ${bilans.applied}.`;
  const nazwy = [...new Set(
    (bilans.skipped ?? [])
      .map((pominiete) => pominiete.lockName ?? pominiete.reason)
      .filter((nazwa) => nazwa !== ''),
  )].slice(0, POWODY_W_KOMUNIKACIE);
  const powod = nazwy.length === 0 ? '' : ` (${nazwy.join(', ')})`;
  return `${czynnosc} — zmieniono ${bilans.applied}, pominięto ${bilans.skippedCount}${powod}.`;
}

export async function wskazDokumentOkna(
  kanal: Kanal,
  idOkna: string,
): Promise<StudioDocument | null> {
  if (idOkna === '') return null;
  const wynik = await wywolaj(kanal, Command.StudioDocumentOpen, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) return null;
  return wynik.wynik.document;
}

// Postać idzie z odczytu rdzenia, bo kanwa niesie sam tekst; odmowa odczytu
// postaci sprowadza zapis do samej treści, zamiast go wstrzymać.
export async function zapiszDokumentZPostacia(
  kanal: Kanal,
  idDokumentu: string,
  tresc: string,
): Promise<StudioDocument | null> {
  const postac = await wywolaj(kanal, Command.StudioDocumentFormGet, { documentId: idDokumentu });
  if (postac.udany && postac.wynik !== undefined) {
    const zapis = await wywolaj(kanal, Command.StudioDocumentFormSave, {
      documentId: idDokumentu,
      form: postac.wynik.form,
      content: tresc,
      createVersion: true,
    });
    if (zapis.udany && zapis.wynik !== undefined) {
      if (zapis.wynik.balance.skippedCount > 0) {
        oglos('Studio', nazwijBilans('Zapis dokumentu', zapis.wynik.balance), 'ostrzezenie');
      }
      return zapis.wynik.document;
    }
    oglos('Studio', nazwijOdmowe('Zapis postaci dokumentu', zapis.blad), 'ostrzezenie');
  }
  const zapis = await wywolaj(kanal, Command.StudioDocumentSave, {
    documentId: idDokumentu,
    content: tresc,
    createVersion: true,
  });
  if (zapis.udany && zapis.wynik !== undefined) return zapis.wynik.document;
  oglos('Studio', nazwijOdmowe('Zapis dokumentu', zapis.blad), 'blad');
  return null;
}

export function zdejmijTrescPrzykladowaFormatu(korzen: ParentNode): boolean {
  const znak = wskazZnakFormatu(korzen);
  if (znak === null) return false;
  znak.textContent = '';
  return true;
}

// Wskazanie znaku przestawia format na kolejny z kontraktu; treść zostaje
// nietknięta, bo jej przełożenie jest osobnym rozstrzygnięciem Operatora.
export function zwiazFormatDokumentu(
  kanal: Kanal,
  idOkna: string,
  korzen: ParentNode,
): Odsubskrybuj | null {
  const znak = wskazZnakFormatu(korzen);
  if (znak === null) return null;
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  let dokument: StudioDocument | null = null;

  const opisz = (): void => {
    znak.textContent = dokument === null ? '' : `format: ${dokument.format}`;
  };

  const przestaw = async (): Promise<void> => {
    if (dokument === null) {
      oglos('Studio', 'Okno nie prowadzi dokumentu — nie ma czemu przestawić formatu.');
      return;
    }
    const wynik = await wywolaj(kanal, Command.StudioDocumentFormatSet, {
      documentId: dokument.id,
      format: nastepnyFormat(dokument.format),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Studio', nazwijOdmowe('Przestawienie formatu', wynik.blad), 'ostrzezenie');
      return;
    }
    dokument = wynik.wynik.document;
    opisz();
  };

  znak.setAttribute('role', 'button');
  znak.tabIndex = 0;
  znak.addEventListener('click', () => {
    void przestaw();
  }, przy);
  znak.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Enter' && zdarzenie.key !== ' ') return;
    zdarzenie.preventDefault();
    void przestaw();
  }, przy);

  const odlacz = zglosUchwyt(EventType.StudioDocumentChanged, (tresc) => {
    if (tresc.document.windowId !== idOkna) return;
    dokument = tresc.document;
    opisz();
  });

  void wskazDokumentOkna(kanal, idOkna).then((otwarty) => {
    dokument = otwarty;
    opisz();
  });

  return () => {
    sterowanie.abort();
    odlacz();
  };
}

function wskazZnakFormatu(korzen: ParentNode): HTMLElement | null {
  const znak = korzen.querySelector('.st-status .st-status-prawa');
  return znak instanceof HTMLElement ? znak : null;
}

function nastepnyFormat(format: StudioDocumentFormat): StudioDocumentFormat {
  const formaty = Object.values(StudioDocumentFormat);
  const miejsce = formaty.indexOf(format);
  return formaty[(miejsce + 1) % formaty.length] ?? format;
}
