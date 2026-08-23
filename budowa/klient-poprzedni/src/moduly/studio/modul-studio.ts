import './studio.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzAparatZrodlo } from './aparat-zrodlo';
import { utworzObiektZrodlo } from './obiekt-zrodlo';
import { utworzOknoIngestOcrPanel } from './okno-ingest-ocr-panel';
import type { ZrodlaPostaciStudia } from './okno-pracy-z-dokumentem';
import { utworzSzablonZrodlo } from './szablon-zrodlo';
import { utworzTabeleZrodlo } from './tabela-zrodlo';
import { utworzZrodloKontroliStudio } from './zrodlo-kontroli-studio';
import { utworzZrodloPostaciStudio } from './zrodlo-postaci-studio';
import { utworzZrodloWstawienStudio } from './zrodlo-wstawien-studio';
import { utworzOknoPetliWykonawczej } from './okno-petli-wykonawczej';
import { utworzOknoPracyZDokumentem } from './okno-pracy-z-dokumentem';
import { utworzOknoRedakcjiDokumentu } from './okno-redakcji-dokumentu';
import { utworzOknoSessionRepository } from './okno-session-repository';
import { utworzOknoToolsPanel } from './okno-tools-panel';
import { utworzOknoWarsztatuDokumentu } from './okno-warsztatu-dokumentu';
import { utworzOsadzenieModulu } from './osadzenie-modulu';
import { utworzZrodloPetli } from './petla-zrodlo';
import { utworzZapleczePrzybornika } from './przybornik-zaplecze';
import { utworzStanStudio } from './stan-studio';
import { utworzZrodloAkcjiStudio } from './zrodlo-akcji-studio';
import { utworzZrodloDokumentuStudio } from './zrodlo-dokumentu-studio';
import { utworzZrodloMaterialuStudio } from './zrodlo-materialu-studio';
import { utworzZrodloOsadzenia } from './zrodlo-osadzenia';
import { utworzZrodloPracyStudio } from './zrodlo-pracy-studio';
import { utworzZrodloPrzekazania } from './zrodlo-przekazania';
import { utworzZrodloWarsztatuDokumentu } from './zrodlo-warsztatu-dokumentu';

/**
 * Moduł Studio — jedno okno pracy z dokumentem i pięć okien bez powierzchni
 * tekstowej.
 *
 * ── Co zeszło się w jedno ───────────────────────────────────────────────────
 * Studio Editor, kanwa tekstowa, Preview Window i Diff/Grep Panel pracowały nad
 * TĄ SAMĄ treścią w czterech miejscach, a przy dwóch dokumentach dawało to sześć
 * okien. Zeszły się w okno pracy z dokumentem: treść z formatowaniem na kartce,
 * podgląd wydania i różnica jako tryby jej widoku, wynik modelu jako zmiana
 * oznaczona w miejscu. Kanwa tekstowa przestała istnieć — wynik operacji
 * kontekstowej wchodzi wprost do dokumentu.
 *
 * ── Układ ───────────────────────────────────────────────────────────────────
 * Ingest/OCR Panel stoi najwyżej, bo cyfryzacja poprzedza redakcję: dopiero jego
 * wynik daje oknu pracy treść, gdy materiałem wejściowym jest skan. Okno pracy
 * jest wiodące i stoi w pasie następnym obok Tools Panelu, który niesie pełny
 * wykaz operacji. Warsztat dokumentu i redakcja pracują na materiale wniesionym
 * do okna, więc stoją niżej. Session Repository jest zarządcą i zamyka układ,
 * bo dotyczy całej sesji, a nie bieżącej czynności.
 *
 * Trzy źródła modułu sięgają poza obszar `studio`, bo obszar ten nie niesie ani
 * cyfryzacji, ani zamiany formatu, ani przygotowania obrazu: dokumenty
 * (`document.*`), materiał (`image.*`, `archive.*`) i przekazanie kontekstu
 * (`context.transfer`). Kontrakt dzieli się po rodzajach czynności, nie po
 * modułach.
 *
 * `stan-studio` jest jeden na cały moduł, więc przywrócenie wersji
 * w repozytorium przestawia treść edytora i podgląd naraz. Odczyty idą
 * równolegle i nie gaszą się nawzajem — odmowa `studio.repository.list` zostaje
 * w Session Repository i nie zabiera Tools Panelowi rejestru akcji.
 */
export interface ModulStudio {
  /** Element osadzany w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Zleca odczyt osadzenia, rejestru akcji i historii wersji. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza subskrypcję zdarzeń rdzenia. */
  rozlacz(): void;
}

export function utworzModulStudio(kanal: Kanal): ModulStudio {
  const stan = utworzStanStudio(kanal);
  const akcje = utworzZrodloAkcjiStudio(kanal);
  const przekazanie = utworzZrodloPrzekazania(kanal);
  const dokumenty = utworzZrodloDokumentuStudio(kanal);
  const material = utworzZrodloMaterialuStudio(kanal);
  const osadzenie = utworzOsadzenieModulu(stan, utworzZrodloOsadzenia(kanal));

  // Zaplecze przybornika stoi na poziomie modułu, bo nastawa trybu wykazu
  // operacji dotyczy dwóch okien naraz: pływaka w oknie pracy i stałego panelu
  // obok. Dwie kopie tej nastawy rozjechałyby się przy pierwszym przełączeniu.
  const przybornikZaplecze = utworzZapleczePrzybornika(kanal);

  // Siedem źródeł postaci dokumentu, kontroli pracy i wstawień. Stoją na
  // poziomie modułu, a nie w oknie pracy, bo narzędziownia cyfryzacji i okno
  // pracy wołają to samo źródło wniesienia: dwa źródła nad jedną rodziną komend
  // byłyby dwiema drogami do jednego rdzenia, a wtedy jedna z nich milczy.
  const zrodlaPostaci: ZrodlaPostaciStudia = {
    postaci: utworzZrodloPostaciStudio(kanal),
    kontrola: utworzZrodloKontroliStudio(kanal),
    wstawienia: utworzZrodloWstawienStudio(kanal),
    tabele: utworzTabeleZrodlo(kanal),
    obiekty: utworzObiektZrodlo(kanal),
    aparat: utworzAparatZrodlo(kanal),
    szablony: utworzSzablonZrodlo(kanal),
  };

  const cyfryzacja = utworzOknoIngestOcrPanel(
    stan,
    dokumenty,
    material,
    zrodlaPostaci.wstawienia,
  );
  const narzedzia = utworzOknoToolsPanel(stan, akcje, {
    naTryb: (tryb) => praca.ustawTrybOperacji(tryb),
  });
  const repozytorium = utworzOknoSessionRepository(stan);
  const zrodloWarsztatu = utworzZrodloWarsztatuDokumentu(kanal);
  const warsztat = utworzOknoWarsztatuDokumentu(stan, zrodloWarsztatu);
  const praca = utworzOknoPracyZDokumentem(
    stan,
    utworzZrodloPracyStudio(kanal),
    akcje,
    dokumenty,
    przekazanie,
    przybornikZaplecze,
    zrodlaPostaci,
    {
      naWiecejOperacji: () => narzedzia.przenieOgnisko(),
      naWarsztat: () => warsztat.przenieOgnisko(),
      naTrybOperacji: (tryb) => narzedzia.ustawTryb(tryb),
    },
  );
  // Redakcja bierze to samo źródło co warsztat: obie drogi wołają rdzeń tak
  // samo — nazwą komendy ze stałych kontraktu i treścią złożoną z formularza.
  const redakcja = utworzOknoRedakcjiDokumentu(stan, zrodloWarsztatu);

  // Okno pętli wykonawczej NIE dostaje stałej kolumny w tym układzie: wchodzi
  // nakładką na żądanie, znacznikiem przebiegu, i schodzi po zwinięciu — bo
  // powierzchnia należy do dokumentu. Stała kolumna jest w nim trybem do wyboru
  // Operatora, przestawianym w samym oknie.
  const petla = utworzOknoPetliWykonawczej(stan, utworzZrodloPetli(kanal), {
    // Kliknięcie zadania prowadzi do MIEJSCA, którego zadanie dotyczy: zakres
    // zadania staje zaznaczeniem dokumentu, a okno pracy wchodzi w pole widzenia.
    // Zakresu nie ma tylko przy zadaniu na całym dokumencie — wtedy zostaje samo
    // przewinięcie, a zaznaczenia nie ruszamy, żeby nie zgubić tego Operatora.
    pokazWynikZadania: (zadanie) => {
      if (zadanie.rangeStart !== undefined && zadanie.rangeEnd !== undefined) {
        stan.ustawZaznaczenie({ poczatek: zadanie.rangeStart, koniec: zadanie.rangeEnd });
      }
      praca.element.scrollIntoView({ block: 'nearest' });
    },
    pokazFragment: (od, do_) => {
      stan.ustawZaznaczenie({ poczatek: od, koniec: do_ });
      praca.element.scrollIntoView({ block: 'nearest' });
    },
  });

  const pasWiodacy = document.createElement('div');
  pasWiodacy.className = 'ms-modul__pas ms-modul__pas--wiodacy';
  pasWiodacy.append(praca.element, narzedzia.element);

  // Pas pętli stoi NAD pasem wiodącym, bo niesie znacznik przebiegu — jedyną
  // drogę do rozwinięcia okna. Sam znacznik to jeden przycisk; okno wchodzi
  // nakładką nad treścią i dokumentowi niczego nie zabiera, dopóki jest zwinięte.
  const pasPetli = document.createElement('div');
  pasPetli.className = 'ms-modul__pas ms-modul__pas--petla';
  pasPetli.append(petla.element);

  // Warsztat dokumentu stoi pod pasem skutków, a nad repozytorium: pracuje na
  // materiale wniesionym do okna, a nie na treści edytora, więc nie należy ani
  // do pasa wiodącego, ani do pasa podglądu.
  const pasWarsztatu = document.createElement('div');
  pasWarsztatu.className = 'ms-modul__pas ms-modul__pas--warsztat';
  pasWarsztatu.append(warsztat.element, redakcja.element);

  const element = document.createElement('div');
  element.className = 'ms-modul';
  element.dataset['modul'] = 'studio';
  element.setAttribute('aria-label', 'Moduł Studio — okna operacyjne');
  element.append(
    osadzenie.element,
    cyfryzacja.element,
    pasPetli,
    pasWiodacy,
    pasWarsztatu,
    repozytorium.element,
  );

  function odswiezWszystkie(): void {
    cyfryzacja.odswiez();
    praca.odswiez();
    narzedzia.odswiez();
    repozytorium.odswiez();
    warsztat.odswiez();
    redakcja.odswiez();
    petla.odswiez();
  }

  const odsubskrybuj = stan.obserwuj(odswiezWszystkie);
  odswiezWszystkie();

  return {
    element,

    async wczytaj(idSesji) {
      // Osadzenie idzie pierwsze, bo dopiero ono ustala `windowId` wymagany
      // przez komendy modułu. Odczyty zależne od dokumentu ruszają po nim.
      await osadzenie.wczytaj(idSesji);
      await Promise.all([
        narzedzia.wczytaj(),
        repozytorium.wczytaj(),
        warsztat.wczytaj(),
        redakcja.wczytaj(),
        praca.wczytaj(),
        petla.wczytaj(),
      ]);
    },

    rozlacz() {
      odsubskrybuj();
      // Tools Panel niesie menu z biblioteki, a to zakłada nasłuch dokumentu —
      // rozbiórka modułu bez jego zwinięcia zostawiłaby nasłuch przy wyjętym
      // z ekranu widoku. Subskrypcja zdarzeń rdzenia i nasłuch menu to dwa
      // różne byty i oba mają tu swoje zamknięcie.
      narzedzia.zamknij();
      // Okno pętli trzyma subskrypcję `studio.chain.progressed` — nasłuch
      // zostawiony przy widoku wyjętym z ekranu odświeżałby okno, którego nie ma.
      petla.zamknij();
      stan.rozlacz();
    },
  };
}
