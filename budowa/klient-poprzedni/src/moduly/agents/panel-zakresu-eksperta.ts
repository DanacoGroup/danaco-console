import type { Agent, AgentPolicy, IsolationTechnicalSwitch } from '../../../../shared/contract';
import { IsolationTechnicalScope } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { StanAgentow } from './stan-agentow';
import type { ZrodloZakresuEksperta } from './zrodlo-zakresu-eksperta';

/**
 * Panel zakresu działania eksperta — trzy grupy Permissions Center, których
 * wiersz przełącznika unieść nie potrafi.
 *
 * Grupa „dostęp do rozszerzeń" i „dostęp do modułów" są wartościami wyliczenia
 * uprawnień i stoją wyżej jako zwykłe wiersze. Tutaj mieszkają trzy rzeczy,
 * które wartością logiczną nie są:
 *
 *   - moduły zastosowania — wykaz kodów, nie przełącznik; wykaz PUSTY znaczy
 *     brak ograniczenia i jest stanem wyjściowym;
 *   - izolacja techniczna — macierz ośmiu suwaków, identyczna z macierzą okna
 *     konfiguracji punktów izolacji, bo macierz jest jedna;
 *   - MultitaskingAI — para „Subagent Network plus górna liczba podagentów".
 *
 * Pod nimi stoi polityka efektywna: podgląd wynikowego zestawu ustawień po
 * uwzględnieniu wszystkich grup. Podgląd, nie bramka — komenda niczego nie
 * zapisuje i niczego nie rozstrzyga.
 *
 * KAŻDE ZAWĘŻENIE STĄD JEST ZAWĘŻENIEM EGZEKWOWANYM. Rdzeń czyta te zapisy
 * przy nakładaniu eksperta na okno i przy powołaniu podagentów, i odmawia.
 * Panel mówi o tym Operatorowi wprost, żeby nie wziął suwaka za ozdobę.
 */
export interface PanelZakresuEksperta {
  element: HTMLElement;
  /** Wczytuje zakres wybranego eksperta; brak wyboru czyści panel. */
  odswiez(): void;
}

/** Osiem zakresów izolacji wraz z ich brzmieniem dla Operatora. */
const OPISY_ZAKRESOW: Record<IsolationTechnicalScope, string> = {
  [IsolationTechnicalScope.WorkingDirectory]: 'katalog roboczy sesji',
  [IsolationTechnicalScope.ProcessEnvironment]: 'środowisko procesu',
  [IsolationTechnicalScope.ModelDataDirectory]: 'katalog danych i konfiguracji modelu',
  [IsolationTechnicalScope.NetworkAccess]: 'dostęp sieciowy',
  [IsolationTechnicalScope.FileAccess]: 'zakres odczytu i zapisu plików',
  [IsolationTechnicalScope.AccountToken]: 'konto i token per sesja',
  [IsolationTechnicalScope.ProcessModel]: 'model procesu',
  [IsolationTechnicalScope.ExecutionServer]: 'serwer wykonania',
};

export function utworzPanelZakresuEksperta(
  stan: StanAgentow,
  zrodlo: ZrodloZakresuEksperta,
): PanelZakresuEksperta {
  const odpowiedz = utworzWierszOdpowiedzi();

  // ── Moduły zastosowania ────────────────────────────────────────────────────
  const moduly = poleTekstowe({
    etykieta: 'Moduły zastosowania',
    podpowiedz: 'kody rozdzielone przecinkiem; puste znaczy wszystkie',
    opis:
      'Moduły, w których ekspert pojawia się jako wykonawca doraźny. WYKAZ PUSTY ZNACZY ' +
      'BRAK OGRANICZENIA i jest stanem wyjściowym. Zawężenie jest egzekwowane: rdzeń ' +
      'odmówi nałożenia tego eksperta na okno modułu spoza wykazu.',
  });
  const zapiszModuly = przycisk('Zapisz moduły', 'dn-btn dn-btn--sm dn-btn--zarys');

  // ── Izolacja techniczna ────────────────────────────────────────────────────
  const macierz = document.createElement('ul');
  macierz.className = 'da-uprawnienia da-uprawnienia--izolacja';

  // ── MultitaskingAI ─────────────────────────────────────────────────────────
  const podagenciCzynne = document.createElement('input');
  podagenciCzynne.type = 'checkbox';
  podagenciCzynne.className = 'dn-przelacznik';

  const granica = document.createElement('input');
  granica.type = 'number';
  granica.className = 'dn-pole-kontrolka';
  granica.min = '1';
  granica.max = '15';

  const zapiszPodagentow = przycisk('Zapisz Subagent Network', 'dn-btn dn-btn--sm dn-btn--zarys');

  const wierszPodagentow = document.createElement('div');
  wierszPodagentow.className = 'dn-pole dm-pole';
  const etykietaPodagentow = document.createElement('span');
  etykietaPodagentow.className = 'da-uprawnienia__etykieta';
  etykietaPodagentow.textContent = 'Subagent Network — górna liczba podagentów (1–15)';
  const opisPodagentow = document.createElement('p');
  opisPodagentow.className = 'dn-pole-opis';
  opisPodagentow.textContent =
    'Wyłączenie usuwa możliwość uruchamiania podagentów przez tego eksperta w dowolnej roli. ' +
    'Rdzeń czyta to przy powołaniu i odmawia — suwak nie jest opisem, tylko regułą.';
  wierszPodagentow.append(etykietaPodagentow, podagenciCzynne, granica, zapiszPodagentow, opisPodagentow);

  // ── Zdjęcie wpisów uprawnień ───────────────────────────────────────────────
  const zdejmij = przycisk('Zdejmij wszystkie zawężenia', 'dn-btn dn-btn--sm dn-btn--zarys');

  // ── Polityka efektywna ─────────────────────────────────────────────────────
  const politykaTresc = document.createElement('pre');
  politykaTresc.className = 'dn-kod da-polityka';
  const odczytajPolityke = przycisk('Przelicz politykę efektywną', 'dn-btn dn-btn--sm dn-btn--zarys');

  const element = document.createElement('section');
  element.className = 'da-panel da-panel--zakres';
  element.dataset['panel'] = 'zakres-eksperta';

  const tytul = document.createElement('h4');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Zakres działania eksperta';

  element.append(
    tytul,
    moduly.element,
    zapiszModuly,
    macierz,
    wierszPodagentow,
    zdejmij,
    odczytajPolityke,
    politykaTresc,
    odpowiedz.element,
  );

  /** Ekspert, którego zakres stoi w panelu — po nim poznajemy zmianę wyboru. */
  let pokazany = '';
  let przelaczniki: IsolationTechnicalSwitch[] = [];

  function ekspertPanelu(): Agent | null {
    return stan.wybrany();
  }

  async function wczytajIzolacje(idEksperta: string): Promise<void> {
    const wynik = await zrodlo.izolacja(idEksperta);
    if (!wynik.udany || wynik.wynik === undefined) {
      macierz.replaceChildren();
      odpowiedz.pokaz(
        opisOdmowy('Odczyt izolacji technicznej', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    przelaczniki = wynik.wynik.switches;
    narysujMacierz();
  }

  function narysujMacierz(): void {
    macierz.replaceChildren(
      ...przelaczniki.map((przelacznik) => {
        const kontrolka = document.createElement('input');
        kontrolka.type = 'checkbox';
        kontrolka.className = 'dn-przelacznik';
        kontrolka.checked = przelacznik.isolated;
        kontrolka.addEventListener('change', () => {
          void zapiszZakres(przelacznik.scope, kontrolka.checked);
        });

        const etykieta = document.createElement('span');
        etykieta.className = 'da-uprawnienia__etykieta';
        etykieta.textContent = OPISY_ZAKRESOW[przelacznik.scope] ?? przelacznik.scope;

        const wiersz = document.createElement('li');
        wiersz.className = 'da-uprawnienia__wiersz';
        wiersz.dataset['zakres'] = przelacznik.scope;
        wiersz.append(kontrolka, etykieta);
        return wiersz;
      }),
    );
  }

  async function zapiszZakres(zakres: IsolationTechnicalScope, odciety: boolean): Promise<void> {
    const ekspert = ekspertPanelu();
    if (ekspert === null) return;
    odpowiedz.pokaz(`Zapis zakresu ${zakres}…`, true);
    const wynik = await zrodlo.zapiszIzolacje(ekspert.id, [{ scope: zakres, isolated: odciety }]);
    if (!wynik.udany || wynik.wynik === undefined) {
      // Suwak wraca do stanu rdzenia: przestawiony przez przeglądarkę nie ma
      // prawa zostać świadectwem zapisu, którego nie było.
      void wczytajIzolacje(ekspert.id);
      odpowiedz.pokaz(
        opisOdmowy('Zapis izolacji technicznej', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    przelaczniki = wynik.wynik.switches;
    narysujMacierz();
    const wpis = przelaczniki.find((kandydat) => kandydat.scope === zakres);
    if (wpis === undefined || wpis.isolated !== odciety) {
      odpowiedz.pokaz(
        `Rdzeń przyjął wywołanie, ale zakres ${zakres} wrócił w innym stanie — zapis się nie odbył.`,
        false,
      );
      return;
    }
    odpowiedz.pokaz(
      `Zakres ${zakres} ${odciety ? 'odcięty' : 'otwarty'}. Ustawienie obowiązuje wszędzie, ` +
        'gdzie ekspert działa, dopóki nie nadpisze go reguła z poziomu bardziej szczegółowego.',
      true,
    );
  }

  async function zapiszModulyEksperta(): Promise<void> {
    const ekspert = ekspertPanelu();
    if (ekspert === null) {
      odpowiedz.pokaz('Wybierz eksperta w Agent Builderze — moduły dotyczą jednego.', false);
      return;
    }
    const kody = moduly.kontrolka.value
      .split(',')
      .map((kod) => kod.trim())
      .filter((kod) => kod !== '');
    odpowiedz.pokaz('Zapis modułów zastosowania…', true);
    const wynik = await zrodlo.ustawModuly(ekspert.id, kody);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Zapis modułów zastosowania', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    stan.wchlon(wynik.wynik.agent);
    const oddane = wynik.wynik.agent.moduleCodes;
    odpowiedz.pokaz(
      oddane.length === 0
        ? 'Ograniczenie zdjęte — ekspert pojawia się jako wykonawca we wszystkich modułach.'
        : `Ekspert pracuje w modułach: ${oddane.join(', ')}. W pozostałych rdzeń odmówi nałożenia.`,
      true,
    );
  }

  async function zapiszSubagentNetwork(): Promise<void> {
    const ekspert = ekspertPanelu();
    if (ekspert === null) {
      odpowiedz.pokaz('Wybierz eksperta w Agent Builderze.', false);
      return;
    }
    const czynny = podagenciCzynne.checked;
    const liczba = Number.parseInt(granica.value, 10);
    odpowiedz.pokaz('Zapis Subagent Network…', true);
    const wynik = await zrodlo.ustawPodagentow(
      ekspert.id,
      czynny,
      czynny && Number.isFinite(liczba) ? liczba : undefined,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      odswiezPola(ekspert);
      odpowiedz.pokaz(
        opisOdmowy('Zapis Subagent Network', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    stan.wchlon(wynik.wynik.agent);
    const zapisana = wynik.wynik.agent.subagentLimit;
    odpowiedz.pokaz(
      zapisana === 0
        ? 'Subagent Network wyłączony — powołanie podagentów przez tego eksperta dostanie odmowę.'
        : `Subagent Network czynny, górna liczba podagentów: ${zapisana}.`,
      true,
    );
  }

  async function zdejmijZawezenia(): Promise<void> {
    const ekspert = ekspertPanelu();
    if (ekspert === null) {
      odpowiedz.pokaz('Wybierz eksperta w Agent Builderze — zdjęcie dotyczy jednego.', false);
      return;
    }
    odpowiedz.pokaz('Zdejmowanie wpisów uprawnień…', true);
    const wynik = await zrodlo.zdejmijUprawnienie({ idEksperta: ekspert.id });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Zdjęcie wpisów uprawnień', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    stan.wchlon({ ...ekspert, permissions: wynik.wynik.permissions });
    odpowiedz.pokaz(
      wynik.wynik.removed === 0
        ? 'Nie było czego zdejmować — ekspert miał już pełny dostęp operacyjny.'
        : `Zdjęto ${wynik.wynik.removed} wpisów. Wiersze zakresu znikają w całości, ` +
            'więc wraca stan „brak ustawienia = wartość domyślna”.',
      true,
    );
  }

  async function przeliczPolityke(): Promise<void> {
    const ekspert = ekspertPanelu();
    if (ekspert === null) {
      odpowiedz.pokaz('Wybierz eksperta w Agent Builderze.', false);
      return;
    }
    const wynik = await zrodlo.polityka(ekspert.id);
    if (!wynik.udany || wynik.wynik === undefined) {
      politykaTresc.textContent = '';
      odpowiedz.pokaz(
        opisOdmowy('Odczyt polityki efektywnej', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    politykaTresc.textContent = opisPolityki(wynik.wynik.policy);
    odpowiedz.pokaz('Polityka efektywna przeliczona — podgląd, nie zapis.', true);
  }

  function odswiezPola(ekspert: Agent): void {
    moduly.kontrolka.value = ekspert.moduleCodes.join(', ');
    podagenciCzynne.checked = ekspert.subagentLimit > 0;
    granica.value = String(ekspert.subagentLimit === 0 ? 15 : ekspert.subagentLimit);
  }

  zapiszModuly.addEventListener('click', () => void zapiszModulyEksperta());
  zapiszPodagentow.addEventListener('click', () => void zapiszSubagentNetwork());
  zdejmij.addEventListener('click', () => void zdejmijZawezenia());
  odczytajPolityke.addEventListener('click', () => void przeliczPolityke());

  return {
    element,

    odswiez() {
      const ekspert = ekspertPanelu();
      if (ekspert === null) {
        pokazany = '';
        macierz.replaceChildren();
        politykaTresc.textContent = '';
        moduly.kontrolka.value = '';
        return;
      }
      if (ekspert.id === pokazany) return;
      pokazany = ekspert.id;
      odpowiedz.wyczysc();
      politykaTresc.textContent = '';
      odswiezPola(ekspert);
      void wczytajIzolacje(ekspert.id);
    },
  };
}

/**
 * Podgląd polityki efektywnej w formie tekstowej.
 *
 * Kształt jest ten sam co w podglądzie okna konfiguracji punktów izolacji:
 * blok techniczny do przeczytania i skopiowania, nie tabela do klikania.
 */
function opisPolityki(polityka: AgentPolicy): string {
  const wiersze: string[] = [];
  wiersze.push(`ekspert: ${polityka.agentId}`);
  wiersze.push(
    `moduly: ${polityka.moduleCodes.length === 0 ? '— bez ograniczenia' : polityka.moduleCodes.join(', ')}`,
  );
  wiersze.push(
    `subagent_network: ${polityka.subagentEnabled ? 'tak' : 'nie'} (granica ${polityka.subagentLimit})`,
  );
  wiersze.push('izolacja_techniczna:');
  for (const przelacznik of polityka.technicalSwitches) {
    wiersze.push(`  ${przelacznik.scope}: ${przelacznik.isolated ? 'odcięty' : 'otwarty'}`);
  }
  if (polityka.permissions.length === 0) {
    wiersze.push('uprawnienia: pełny dostęp operacyjny — ani jednego zawężenia');
  } else {
    wiersze.push('uprawnienia:');
    for (const wpis of polityka.permissions) {
      const zakres = wpis.scope === undefined || wpis.scope === '' ? 'cała grupa' : wpis.scope;
      wiersze.push(`  ${wpis.group} / ${zakres}: ${wpis.granted ? 'przyznane' : 'odebrane'}`);
    }
  }
  if (polityka.overriddenBy !== undefined) {
    wiersze.push(`nadpisane na poziomie zasiegu: ${polityka.overriddenBy}`);
  }
  return wiersze.join('\n');
}
