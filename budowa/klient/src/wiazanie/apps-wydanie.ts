// Czynności wydania w oknie Apps: pakiet wraz z manifestem, podpisem i
// ogłoszeniem, wdrożenia, domena, skala, zdrowie, dzienniki i zmienne.
import {
  AppDeployEnvironment,
  AppPackageFormat,
  Command,
  ExtensionKind,
} from '../../../shared/contract.ts';
import type { AppEnvironmentVariable, AppValidationIssue } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { potwierdzone, wypelnij } from './apps-wytworzenie.ts';

const NAGLOWEK = 'Wydanie aplikacji';

/* Oznaczenie pakietu żyje w oknie, bo kontrakt nie ma wykazu pakietów: bierze
   się z budowy albo zapisu manifestu i służy podpisowi, kontroli i ogłoszeniu. */
const PAKIETY = new Map<string, string>();

const PRZYCISKI: ReadonlyArray<readonly [string, string]> = [
  ['manifest', 'Zapisz manifest'],
  ['pakiet', 'Zbuduj pakiet'],
  ['pakiet-sprawdz', 'Sprawdź pakiet'],
  ['pakiet-podpisz', 'Podpisz pakiet'],
  ['pakiet-oglos', 'Ogłoś pakiet'],
  ['wdroz', 'Wdroż na środowisko'],
  ['domena', 'Ustaw domenę'],
  ['skala', 'Ustaw skalę'],
  ['zdrowie', 'Odczytaj zdrowie'],
  ['dziennik-wdrozenia', 'Dziennik wdrożenia'],
  ['dziennik-uslugi', 'Dziennik usługi'],
  ['zmienne', 'Odczytaj zmienne'],
  ['zmienna-ustaw', 'Ustaw zmienną'],
];

/*
zwiazWydanieAplikacji stawia pas nad panelem wdrożeń.

Pas stoi nad ciałem panelu, bo ciało jest wymieniane przy każdym odświeżeniu
wykazu wdrożeń; pas postawiony w nim znikałby razem z pozycjami.
*/
export function zwiazWydanieAplikacji(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  przy: AddEventListenerOptions,
): void {
  postawPas(korzen);
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-wydanie]')?.dataset.wydanie;
    if (czynnosc === undefined) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, korzen, czynnosc, idOkna(), odswiez);
  }, przy);
}

function postawPas(korzen: Element): void {
  const panel = korzen.querySelector('#panel-deployment');
  const cialo = panel?.querySelector('.sta-okno-tresc');
  if (panel === null || panel === undefined || cialo === null || cialo === undefined) return;
  if (panel.querySelector('[data-wydanie-wpis]') !== null) return;
  const pas = korzen.ownerDocument.createElement('div');
  pas.className = 'dn-pas-dzialan';
  const pole = korzen.ownerDocument.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole dn-pole--sm';
  pole.placeholder = 'Środowisko, nazwa albo wartość';
  pole.setAttribute('aria-label', 'Środowisko, nazwa albo wartość czynności wydania');
  pole.dataset.wydanieWpis = '';
  pas.appendChild(pole);
  for (const [czynnosc, etykieta] of PRZYCISKI) {
    const wezel = korzen.ownerDocument.createElement('button');
    wezel.type = 'button';
    wezel.className = 'dn-btn dn-btn--duch dn-btn--sm';
    wezel.dataset.wydanie = czynnosc;
    wezel.textContent = etykieta;
    pas.appendChild(wezel);
  }
  panel.insertBefore(pas, cialo);
}

function pole(korzen: Element): string {
  const wezel = korzen.querySelector<HTMLInputElement>('[data-wydanie-wpis]');
  return (wezel?.value ?? '').trim();
}

function czesci(korzen: Element): string[] {
  return pole(korzen).split('|').map((czesc) => czesc.trim());
}

/* Bez wskazania w polu czynność idzie na środowisko robocze: wdrożenie na
   produkcję ma być wyborem wpisanym wprost, nie skutkiem pustego pola. */
function srodowisko(oznaczenie: string): AppDeployEnvironment {
  if (oznaczenie === 'production') return AppDeployEnvironment.Production;
  if (oznaczenie === 'staging') return AppDeployEnvironment.Staging;
  return AppDeployEnvironment.Dev;
}

function pakiet(idOkna: string): string {
  const zapamietany = PAKIETY.get(idOkna) ?? '';
  if (zapamietany !== '') return zapamietany;
  oglos(NAGLOWEK, 'Okno nie zna jeszcze pakietu — zbuduj go albo zapisz manifest.', 'ostrzezenie');
  return '';
}

async function wykonaj(
  kanal: Kanal,
  korzen: Element,
  czynnosc: string,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna aplikacji dla tej karty.', 'ostrzezenie');
    return;
  }
  if (czynnosc === 'manifest') return zapiszManifest(kanal, korzen, idOkna);
  if (czynnosc === 'pakiet') return zbudujPakiet(kanal, idOkna);
  if (czynnosc === 'pakiet-sprawdz') return sprawdzPakiet(kanal, korzen, idOkna);
  if (czynnosc === 'pakiet-podpisz') return podpiszPakiet(kanal, korzen, idOkna);
  if (czynnosc === 'pakiet-oglos') return oglosPakiet(kanal, korzen, idOkna);
  if (czynnosc === 'wdroz') return wdroz(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'domena') return ustawDomene(kanal, korzen, idOkna);
  if (czynnosc === 'skala') return ustawSkale(kanal, korzen, idOkna);
  if (czynnosc === 'zdrowie') return odczytajZdrowie(kanal, korzen, idOkna);
  if (czynnosc === 'dziennik-wdrozenia') return odczytajDziennikWdrozenia(kanal, korzen, idOkna);
  if (czynnosc === 'dziennik-uslugi') return odczytajDziennikUslugi(kanal, korzen, idOkna);
  if (czynnosc === 'zmienne') return odczytajZmienne(kanal, korzen, idOkna);
  if (czynnosc === 'zmienna-ustaw') return ustawZmienna(kanal, korzen, idOkna);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: przycisk bez gałęzi
     wyglądałby jak działający. */
  oglos(NAGLOWEK, `Czynność „${czynnosc}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

/* Manifest bierze z pola oznaczenie, nazwę i wersję; rodzajem jest wtyczka,
   bo tak wydaje się aplikację zbudowaną w tym oknie. */
async function zapiszManifest(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [oznaczenie = '', nazwa = '', wersja = ''] = czesci(korzen);
  if (oznaczenie === '' || nazwa === '' || wersja === '') {
    oglos(NAGLOWEK,
      'Manifest potrzebuje oznaczenia, nazwy i wersji, na przykład „sklep | Sklep | 1.0.0".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsPackageManifestSave, {
    windowId: idOkna,
    manifest: { identifier: oznaczenie, name: nazwa, version: wersja, kind: ExtensionKind.Plugin },
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu manifestu.', 'ostrzezenie');
    return;
  }
  PAKIETY.set(idOkna, wynik.wynik.package.id);
  oglos(NAGLOWEK, `Manifest zapisany; pakiet ${wynik.wynik.package.id}.`);
}

async function zbudujPakiet(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsPackageBuild, {
    windowId: idOkna,
    format: AppPackageFormat.Zip,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zbudowania pakietu.', 'ostrzezenie');
    return;
  }
  PAKIETY.set(idOkna, wynik.wynik.package.id);
  oglos(NAGLOWEK, `Pakiet ${wynik.wynik.package.id} zbudowany jako wytwór `
    + `${wynik.wynik.package.artifactRef ?? 'bez oznaczenia'}.`);
}

async function sprawdzPakiet(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cel = pakiet(idOkna);
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.AppsPackageValidate, {
    windowId: idOkna,
    packageId: cel,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia pakietu.', 'ostrzezenie');
    return;
  }
  const uwagi = wynik.wynik.issues;
  if (uwagi.length === 0) {
    oglos(NAGLOWEK, 'Pakiet bez uwag.');
    return;
  }
  wypelnij(korzen, 'panel-deployment',
    uwagi.map((uwaga: AppValidationIssue) => `${uwaga.severity} · ${uwaga.message}`));
  oglos(NAGLOWEK, `Pakiet ma ${uwagi.length} uwag.`, 'ostrzezenie');
}

/* Podpis idzie odwołaniem do klucza w sejfie rdzenia; materiał klucza nie
   przechodzi przez okno ani przez ogłoszenie. */
async function podpiszPakiet(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const odwolanie = pole(korzen);
  if (odwolanie === '') {
    oglos(NAGLOWEK, 'Podpis potrzebuje odwołania do klucza w sejfie rdzenia.', 'ostrzezenie');
    return;
  }
  const cel = pakiet(idOkna);
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.AppsPackageSign, {
    windowId: idOkna,
    packageId: cel,
    signingKeyRef: odwolanie,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił podpisania pakietu.', 'ostrzezenie');
    return;
  }
  const podpis = wynik.wynik.signature;
  oglos(NAGLOWEK, podpis.signed
    ? `Pakiet podpisany na poziomie zaufania ${podpis.trustLevel}.`
    : `Rdzeń nie podpisał pakietu: ${podpis.detail ?? 'bez wyjaśnienia'}.`,
  podpis.signed ? 'informacja' : 'ostrzezenie');
}

/* Ogłoszenie pakietu wystawia go innym kontom i cofnąć się go nie da, więc
   pierwsze naciśnięcie uzbraja przycisk. */
async function oglosPakiet(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cel = pakiet(idOkna);
  if (cel === '') return;
  if (!potwierdzone(korzen, 'data-wydanie', 'pakiet-oglos')) return;
  const wynik = await wywolaj(kanal, Command.AppsPackagePublish, {
    windowId: idOkna,
    packageId: cel,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ogłoszenia pakietu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pakiet ogłoszony jako rozszerzenie „${wynik.wynik.extension.name}".`);
}

/* Wdrożenie na produkcję jest nieodwracalne dla odbiorców usługi, więc
   wymaga drugiego naciśnięcia; środowiska robocze idą od razu. */
async function wdroz(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const cel = srodowisko(pole(korzen));
  if (cel === AppDeployEnvironment.Production && !potwierdzone(korzen, 'data-wydanie', 'wdroz')) return;
  const wynik = await wywolaj(kanal, Command.AppsDeploymentRun, {
    windowId: idOkna,
    environment: cel,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wdrożenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Wdrożenie na ${cel} w stanie ${wynik.wynik.deployment.status}.`);
  odswiez();
}

async function ustawDomene(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [oznaczenie = '', domena = ''] = czesci(korzen);
  if (domena === '') {
    oglos(NAGLOWEK,
      'Domena potrzebuje środowiska i adresu, na przykład „production | sklep.example.com".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsDeploymentDomainSet, {
    windowId: idOkna,
    environment: srodowisko(oznaczenie),
    domain: domena,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia domeny.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Domena ${domena} ustawiona na środowisku ${srodowisko(oznaczenie)}.`);
}

async function ustawSkale(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [oznaczenie = '', liczba = ''] = czesci(korzen);
  const wystapienia = Number.parseInt(liczba, 10);
  if (!Number.isFinite(wystapienia) || wystapienia < 1) {
    oglos(NAGLOWEK,
      'Skala potrzebuje środowiska i liczby wystąpień, na przykład „staging | 3".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsDeploymentScaleSet, {
    windowId: idOkna,
    environment: srodowisko(oznaczenie),
    instances: wystapienia,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia skali.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Skala ${wystapienia} ustawiona na środowisku ${srodowisko(oznaczenie)}.`);
}

async function odczytajZdrowie(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsDeploymentHealthGet, {
    windowId: idOkna,
    environment: srodowisko(pole(korzen)),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu zdrowia.', 'ostrzezenie');
    return;
  }
  const zdrowie = wynik.wynik.health;
  oglos(NAGLOWEK, `Środowisko ${zdrowie.environment}: `
    + `${zdrowie.available ? 'dostępne' : 'niedostępne'}, `
    + `dostępność ${zdrowie.availabilityPercent ?? 'nieznana'}.`,
  zdrowie.available ? 'informacja' : 'ostrzezenie');
}

/* Dziennik czytany jest z wdrożenia stojącego w wykazie najwyżej: okno nie
   prowadzi wskazania wdrożenia, a wykaz podaje ich kolejność. */
async function odczytajDziennikWdrozenia(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.AppsDeploymentList, { windowId: idOkna });
  const cel = wykaz.wynik?.deployments[0]?.id ?? '';
  if (cel === '') {
    oglos(NAGLOWEK, 'To okno nie ma jeszcze żadnego wdrożenia.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsDeploymentLogRead, {
    windowId: idOkna,
    deploymentId: cel,
    limit: 100,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziennika wdrożenia.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-deployment', wynik.wynik.lines);
  if (wynik.wynik.lines.length === 0) oglos(NAGLOWEK, 'Dziennik tego wdrożenia jest pusty.');
}

async function odczytajDziennikUslugi(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsServiceLogRead, { windowId: idOkna, limit: 100 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziennika usługi.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-backend', wynik.wynik.lines);
  if (wynik.wynik.lines.length === 0) oglos(NAGLOWEK, 'Dziennik usługi jest pusty.');
}

/* Wartości tajne rdzeń oddaje odwołaniem do sejfu, nie treścią; wykaz nazywa
   odwołanie, żeby nie udawał, że zna wartość. */
async function odczytajZmienne(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsEnvironmentVariableList, {
    windowId: idOkna,
    environment: srodowisko(pole(korzen)),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu zmiennych.', 'ostrzezenie');
    return;
  }
  const zmienne = wynik.wynik.variables;
  wypelnij(korzen, 'panel-terminal', zmienne.map((zmienna: AppEnvironmentVariable) =>
    `${zmienna.name} = ${zmienna.value ?? `sejf ${zmienna.secretRef ?? 'bez odwołania'}`}`));
  if (zmienne.length === 0) oglos(NAGLOWEK, 'To środowisko nie ma jeszcze zmiennych.');
}

async function ustawZmienna(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [oznaczenie = '', nazwa = '', wartosc = ''] = czesci(korzen);
  if (nazwa === '') {
    oglos(NAGLOWEK,
      'Zmienna potrzebuje środowiska, nazwy i wartości, na przykład „dev | PORT | 8080".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsEnvironmentVariableSet, {
    windowId: idOkna,
    environment: srodowisko(oznaczenie),
    name: nazwa,
    value: wartosc,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia zmiennej.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zmienna ${nazwa} ustawiona na środowisku ${srodowisko(oznaczenie)}.`);
}
