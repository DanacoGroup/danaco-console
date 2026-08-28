import { DebugStepKind, type DebugSession } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { pole, poleLiczbowe, przyciskAkcji, wiersz } from '../../modele/kontrolki-formularza-braki';
import type { StanDevelopera } from './stan-developer';
import { utworzStanTresci } from './stany-okna';
import { cialoZakladki, objasnienieZakladki, pasekZakladki } from './zakladki-okna';
import { rysujZaleznosci, zaleznosci } from './zaleznosci-zewnetrzne';
import type { ZrodloWarsztatu } from './zrodlo-warsztatu';

/**
 * Run & Debug — druga część kolumny monitora modułu Developer, dla sesji
 * debugowania w toku i po jej zakończeniu.
 */
export interface PanelRunDebug {
  element: HTMLElement;
  odswiez(): void;
  zamknij(): void;
}

/** Buduje ciało zakładki Run & Debug kolumny monitora modułu Developer w oknie sesji rdzenia aplikacji. */
export function utworzPanelRunDebug(
  zrodlo: ZrodloWarsztatu,
  stan: StanDevelopera,
): PanelRunDebug {
  const tresc = utworzStanTresci();
  const program = pole('Program do debugowania', 'katalog albo plik pakietu głównego');
  const parametry = pole('Parametry uruchomienia', 'rozdzielone spacją');
  const plikPunktu = pole('Plik punktu przerwania', 'sciezka/do/pliku.go');
  const wierszPunktu = poleLiczbowe('Wiersz punktu przerwania', '42');
  const warunek = pole('Warunek punktu', 'wynik > 3');
  const wyrazenie = pole('Wyrażenie do obliczenia', 'suma');

  let sesja: DebugSession | null = null;
  let ramkaBiezaca = '';

  const postaw = przyciskAkcji('Ustaw punkt przerwania');
  const zdejmij = przyciskAkcji('Zdejmij punkt przerwania');
  const rozpocznij = przyciskAkcji('Rozpocznij sesję debugowania');
  const kontynuuj = przyciskAkcji('Kontynuuj');
  const przejdz = przyciskAkcji('Przejdź (krok)');
  const wejdz = przyciskAkcji('Wejdź');
  const wyjdz = przyciskAkcji('Wyjdź');
  const zatrzymaj = przyciskAkcji('Zakończ sesję');
  const odczytaj = przyciskAkcji('Pokaż stos i zmienne');
  const policz = przyciskAkcji('Wykonaj wyrażenie');

  postaw.addEventListener('click', () => void ustawPunkt(false));
  zdejmij.addEventListener('click', () => void ustawPunkt(true));
  rozpocznij.addEventListener('click', () => void rozpocznijSesje());
  kontynuuj.addEventListener('click', () => void krok(DebugStepKind.Continue));
  przejdz.addEventListener('click', () => void krok(DebugStepKind.StepOver));
  wejdz.addEventListener('click', () => void krok(DebugStepKind.StepInto));
  wyjdz.addEventListener('click', () => void krok(DebugStepKind.StepOut));
  zatrzymaj.addEventListener('click', () => void krok(DebugStepKind.Stop));
  odczytaj.addEventListener('click', () => void pokazZakres());
  policz.addEventListener('click', () => void obliczWyrazenie());

  async function ustawPunkt(zdjecie: boolean): Promise<void> {
    const sciezka = plikPunktu.value.trim() === '' ? stan.sciezka() : plikPunktu.value.trim();
    const numer = Number.parseInt(wierszPunktu.value, 10);
    if (sciezka === '' || Number.isNaN(numer) || numer < 1) {
      tresc.blad('Punkt przerwania wymaga pliku oraz wiersza liczonego od jedynki.');
      return;
    }
    const zadanie: Parameters<ZrodloWarsztatu['punktPrzerwania']>[0] = {
      windowId: stan.okno(),
      path: sciezka,
      line: numer,
    };
    if (zdjecie) zadanie.remove = true;
    else if (warunek.value.trim() !== '') zadanie.condition = warunek.value.trim();

    const wynik = await zrodlo.punktPrzerwania(zadanie);
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Punkt przerwania', wynik.blad), wynik.blad);
      return;
    }
    if (wynik.wynik.breakpoints.length === 0) {
      tresc.pusto('Okno nie ma ani jednego punktu przerwania.');
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      ...wynik.wynik.breakpoints.map((punkt) =>
        akapitPanelu(
          `${punkt.path}:${String(punkt.line)} — ${punkt.kind}` +
            `${punkt.condition === undefined ? '' : ` gdy ${punkt.condition}`}` +
            `${punkt.verified ? '' : ' (niepotwierdzony przez adapter)'}`,
        ),
      ),
    );
  }

  async function rozpocznijSesje(): Promise<void> {
    tresc.ladowanie('Uruchamianie programu pod adapterem debugowania…');
    const zadanie: Parameters<ZrodloWarsztatu['rozpocznijDebugowanie']>[0] = {
      windowId: stan.okno(),
    };
    if (program.value.trim() !== '') zadanie.program = program.value.trim();
    const argumenty = parametry.value.trim().split(/\s+/).filter((czesc) => czesc !== '');
    if (argumenty.length > 0) zadanie.arguments = argumenty;

    const wynik = await zrodlo.rozpocznijDebugowanie(zadanie);
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Rozpoczęcie sesji debugowania', wynik.blad), wynik.blad);
      return;
    }
    sesja = wynik.wynik;
    tresc.potwierdzenie(
      `Sesja ${sesja.id} na adapterze ${sesja.adapter} — stan ${sesja.status}.`,
      true,
    );
  }

  async function krok(rodzaj: DebugStepKind): Promise<void> {
    if (sesja === null) {
      tresc.blad('Kroku nie ma czym wykonać — sesja debugowania nie została rozpoczęta.');
      return;
    }
    const wynik = await zrodlo.sterujDebugowaniem({ sessionId: sesja.id, step: rodzaj });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Sterowanie sesją debugowania', wynik.blad), wynik.blad);
      return;
    }
    sesja = wynik.wynik;
    if (rodzaj === DebugStepKind.Stop) {
      sesja = null;
      ramkaBiezaca = '';
      tresc.potwierdzenie('Sesja debugowania zakończona.', true);
      return;
    }
    // Stan po kroku odczytuje się z rdzenia, a nie z naciśniętego przycisku programu.
    await pokazZakres();
  }

  async function pokazZakres(): Promise<void> {
    if (sesja === null) {
      tresc.blad('Stosu wywołań nie ma skąd wziąć — sesja debugowania nie została rozpoczęta.');
      return;
    }
    const wynik = await zrodlo.zakresDebugowania({ sessionId: sesja.id });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Odczyt stosu i zmiennych', wynik.blad), wynik.blad);
      return;
    }
    const zakres = wynik.wynik;
    if (zakres.frames.length === 0) {
      tresc.pusto(
        'Program biegnie — stos wywołań istnieje wtedy, gdy jest zatrzymany na punkcie przerwania.',
      );
      return;
    }
    ramkaBiezaca = zakres.frames[0]?.id ?? '';
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      akapitPanelu('Stos wywołań:'),
      ...zakres.frames.map((ramka) =>
        akapitPanelu(
          `${ramka.name}${ramka.path === undefined ? '' : ` — ${ramka.path}`}` +
            `${ramka.line === undefined ? '' : `:${String(ramka.line)}`}`,
        ),
      ),
      akapitPanelu('Zmienne:'),
      ...zakres.variables.map((zmienna) =>
        akapitPanelu(
          `${zmienna.name} = ${zmienna.value}` +
            `${zmienna.type === undefined ? '' : ` (${zmienna.type})`}`,
        ),
      ),
    );
  }

  async function obliczWyrazenie(): Promise<void> {
    if (sesja === null || ramkaBiezaca === '') {
      tresc.blad(
        'Wyrażenie liczy się w kontekście ramki stosu — najpierw zatrzymaj program i odczytaj stos.',
      );
      return;
    }
    if (wyrazenie.value.trim() === '') {
      tresc.blad('Obliczenie wymaga wyrażenia.');
      return;
    }
    const wynik = await zrodlo.obliczWyrazenie({
      sessionId: sesja.id,
      frameId: ramkaBiezaca,
      expression: wyrazenie.value.trim(),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Wykonanie wyrażenia', wynik.blad), wynik.blad);
      return;
    }
    tresc.potwierdzenie(
      `${wyrazenie.value.trim()} = ${wynik.wynik.value}` +
        `${wynik.wynik.type === undefined ? '' : ` (${wynik.wynik.type})`}`,
      true,
    );
  }

  const element = cialoZakladki(
    objasnienieZakladki(
      'Punkt przerwania należy do okna, nie do sesji: postawiony raz zostaje na drugi i trzeci ' +
        'bieg, bo rdzeń trzyma go w bazie i podaje adapterowi przy starcie sesji.',
    ),
    objasnienieZakladki(
      'Stan po kroku odczytuje się z rdzenia, a nie z naciśniętego przycisku. Zatrzymanie zgłasza ' +
        'debugowany proces, więc panel po każdym kroku pyta o stos wywołań — rysowanie stanu ' +
        'z samego przycisku pokazywałoby stan życzeniowy.',
    ),
    objasnienieZakladki(
      'Konfiguracje uruchomień to nie to samo co zadania budowania. Zadanie budowania jedzie polem ' +
        '„Zadanie” zakładki Build Output; konfiguracja uruchomienia wskazuje ponadto program ' +
        'i adapter debugowania właściwy językowi.',
    ),
    rysujZaleznosci(zaleznosci(['adapter-dap'])),
    wiersz('Program', program, { klasa: 'dn-pole' }),
    wiersz('Parametry', parametry, { klasa: 'dn-pole' }),
    wiersz('Plik punktu przerwania', plikPunktu, {
      klasa: 'dn-pole',
      objasnienie: 'Puste znaczy plik wskazany w Project Tree.',
    }),
    wiersz('Wiersz punktu przerwania', wierszPunktu, { klasa: 'dn-pole' }),
    wiersz('Warunek punktu', warunek, {
      klasa: 'dn-pole',
      objasnienie: 'Podany warunek czyni punkt warunkowym.',
    }),
    wiersz('Wyrażenie', wyrazenie, { klasa: 'dn-pole' }),
    pasekZakladki(postaw, zdejmij, rozpocznij),
    pasekZakladki(kontynuuj, przejdz, wejdz, wyjdz, zatrzymaj, odczytaj, policz),
    tresc.element,
  );

  return {
    element,
    // Panel nie odpytuje rdzenia sam z siebie: sesja debugowania jest czynnością operatora.
    odswiez: () => undefined,
    zamknij: () => {
      sesja = null;
    },
  };
}

/** akapitPanelu składa jeden wiersz treści panelu Run & Debug kolumny monitora modułu Developer w oknie. */
function akapitPanelu(zdanie: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'mdev-wiersz';
  element.textContent = zdanie;
  return element;
}
