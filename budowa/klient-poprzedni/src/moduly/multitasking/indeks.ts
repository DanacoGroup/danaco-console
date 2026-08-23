import type { Kanal } from '../../protokol/kanal';
import { utworzZrodloSekcjiPaneli } from '../../powloka/zrodlo-sekcji-paneli';
import { widokZSesji, type OpisModulu } from '../rejestracja';
import { utworzWykazKomendRdzenia } from './braki-kontraktu';
import { utworzOknoAnalityka } from './okno-results-analyzer';
import { utworzPanelObsady } from './panel-obsady';
import { utworzPanelZadanWTle } from './panel-zadan-w-tle';
import { utworzZrodloPodagentow } from './zrodlo-podagentow';
import { utworzOknoKoordynatora } from './okno-coordinator-chat';
import { utworzOknoWykonawcy } from './okno-executor-chat';
import { utworzStanMultitaskingu } from './stan-multitaskingu';
import { utworzZrodloBiegu } from './zrodlo-biegu';
import { utworzZrodloOkien } from './zrodlo-okien';

/**
 * Moduł MultitaskingAI — okna ról pętli koordynator–wykonawca.
 *
 * Cztery wystąpienia, trzy typy: Coordinator Chat planuje i steruje, dwa
 * Executor Chat wykonują równolegle, Results Analyzer zestawia wyniki.
 * Koordynator nie tworzy produktu końcowego, a wykonawca nie zarządza procesem.
 *
 * Wszystkie cztery stoją na jednym stanie i jednym źródle biegu — inaczej
 * koordynator widziałby inny przebieg niż wykonawca, którym steruje.
 */

/** Ile okien wykonawców stoi na scenie; wykaz mówi o dwóch (Executor 1 i 2). */
const LICZBA_WYKONAWCOW = 2;

/**
 * Składa scenę czterech okien ról dla wskazanej sesji.
 *
 * Sesja jest tu wymagana, nie opcjonalna: role przypisuje się do OKIEN sesji,
 * więc bez niej nie ma czego obsadzić. Powłoka podaje ją przy wczytaniu.
 */
export function zamontujMultitasking(gospodarz: HTMLElement, kanal: Kanal, sesja: string): {
  zamknij(): void;
} {
  const bieg = utworzZrodloBiegu(kanal);
  const okna = utworzZrodloOkien(kanal);
  const podagenci = utworzZrodloPodagentow(kanal);
  const stan = utworzStanMultitaskingu(bieg, { sesja });
  // Wykaz komend rdzenia powstaje raz na scenę i jest wspólny czterem oknom:
  // cztery osobne pytania o to samo dałyby cztery odpowiedzi do rozjechania.
  const komendy = utworzWykazKomendRdzenia(kanal);

  const scena = document.createElement('div');
  scena.className = 'dm-modul';
  scena.dataset['modul'] = 'multitasking';

  // Panel obsady stoi przed oknami ról, bo bez przydziału ról nie ma czego
  // otworzyć: koordynator i wykonawcy są rolami okien sesji, a nie osobnymi
  // bytami. Panel jest jedynym miejscem, w którym te role się nadaje, a jego
  // meldunki idą do jego własnego stanu treści — scena ich nie przechwytuje,
  // więc odmowa rdzenia nie ma jak zniknąć po drodze.
  // Układ sekcji paneli bierze się z powłoki, nie z modułu: `panel.sections.*`
  // dotyczy okna, a nie dziedziny MultitaskingAI. Moduł tworzy źródło na swoim
  // kanale i podaje je panelowi — mechanika sekcji zostaje jedna.
  const sekcje = utworzZrodloSekcjiPaneli(kanal);
  const obsada = utworzPanelObsady({ zrodlo: okna, stan, sekcje });

  const koordynator = utworzOknoKoordynatora({ okna, bieg, stan, komendy });
  // Wykonawcy dostają to samo źródło podagentów co panel zadań w tle —
  // drugie wywołanie tych samych komend w drugim źródle byłoby kopią.
  const wykonawcy = Array.from({ length: LICZBA_WYKONAWCOW }, (_, i) =>
    utworzOknoWykonawcy({ okna, bieg, stan, podagenci, komendy, numer: i + 1 }),
  );
  const analityk = utworzOknoAnalityka({ okna, bieg, stan, komendy });

  // Panel zadań w tle stoi pod oknami ról, bo mówi o pracy, którą one już
  // zleciły: podagenci należą do okien wykonawców, a nie do osobnej sceny.
  // Jest jeden na całą kartę sesji — `subagent.list` przyjmuje `sessionId`,
  // więc dwa panele po jednym na wykonawcę pytałyby dwa razy o to samo.
  const zadaniaWTle = utworzPanelZadanWTle({ podagenci, bieg, stan });

  scena.append(
    obsada.element,
    koordynator.element,
    ...wykonawcy.map((w) => w.element),
    analityk.element,
    zadaniaWTle.element,
  );
  gospodarz.append(scena);

  // Pytanie o wykaz komend idzie przed pierwszym odświeżeniem okien: dopóki
  // odpowiedź nie wróci, kontrolki bez pokrycia mówią, że pytanie jest w drodze,
  // a nie że rdzeń czegoś nie ma.
  void komendy.odczytaj();
  obsada.odswiez();
  koordynator.odswiez();
  for (const wykonawca of wykonawcy) wykonawca.odswiez();
  analityk.odswiez();
  zadaniaWTle.odswiez();

  return {
    zamknij() {
      // Nasłuchy odpina się przed zdjęciem sceny: zdarzenie, które przyszłoby
      // po `remove()`, odświeżałoby widok już nieistniejący.
      for (const wykonawca of wykonawcy) wykonawca.rozlacz();
      scena.remove();
    },
  };
}

/**
 * Samoopisujący się moduł dla rejestru powłoki.
 *
 * Kod `multitaskingai` odpowiada wierszowi środowiska w katalogu rdzenia —
 * nie literałowi wymyślonemu w kliencie.
 */
export const MODUL: OpisModulu = {
  kod: 'multitaskingai',
  utworzWidok: (kanal) =>
    widokZSesji(kanal, (gospodarz: HTMLElement, k: Kanal, sesja: string) =>
      zamontujMultitasking(gospodarz, k, sesja)),
};
