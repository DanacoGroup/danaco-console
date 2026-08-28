import type {
  IsolationProfile,
  IsolationProfileSaveRequest,
  IsolationSwitch,
  IsolationTechnicalSwitch,
} from '../../../shared/contract';
import { pole, poleTresci, przelacznik, wiersz } from '../modele/kontrolki-formularza-braki';
import {
  POZYCJE_KONTEKSTU,
  POZYCJE_TECHNICZNE,
  etykietaKontekstu,
  etykietaTechniczna,
} from './katalog-izolacji';

/** Formularz profilu izolacji — nazwa, opis i komplet jedenastu przełączników w jednym zestawie zapisu. */
export interface FormularzProfilu {
  /** Element montowany w obszarze Profile. */
  element: HTMLElement;
  /** Treść żądania `isolation.profile.save` złożona ze stanu formularza. */
  zadanie(): IsolationProfileSaveRequest;
  /** Wypełnia formularz profilem wczytanym z rdzenia i przestawia zapis na zmianę tego profilu. */
  wypelnij(profil: IsolationProfile): void;
  /** Czyści formularz i przestawia zapis na zakładanie nowego profilu. */
  wyczysc(): void;
}

/** Przełącznik z etykietą wartości aktualizowaną przy każdym przestawieniu tej samej kontrolki formularza. */
interface Przelacznik {
  element: HTMLElement;
  kontrolka: HTMLInputElement;
}

export function utworzFormularzProfilu(): FormularzProfilu {
  let profileId: string | undefined;

  const nazwa = pole('Nazwa profilu', 'np. Ostry rozdział — projekt klienta');
  const opis = poleTresci('Opis profilu', 2, 'Do czego ten zestaw służy i kiedy go przypisujesz.');

  const przelacznikiKontekstu = new Map<string, Przelacznik>();
  const przelacznikiTechniczne = new Map<string, Przelacznik>();

  const grupaKontekstu = document.createElement('div');
  grupaKontekstu.className = 'pi-profile__przelaczniki';
  for (const pozycja of POZYCJE_KONTEKSTU) {
    const zbudowany = zbudujPrzelacznik(pozycja.nazwa, etykietaKontekstu);
    przelacznikiKontekstu.set(pozycja.kind, zbudowany);
    grupaKontekstu.append(zbudowany.element);
  }

  const grupaTechniczna = document.createElement('div');
  grupaTechniczna.className = 'pi-profile__przelaczniki';
  for (const pozycja of POZYCJE_TECHNICZNE) {
    const zbudowany = zbudujPrzelacznik(pozycja.nazwa, etykietaTechniczna);
    przelacznikiTechniczne.set(pozycja.scope, zbudowany);
    grupaTechniczna.append(zbudowany.element);
  }

  const cel = document.createElement('p');
  cel.className = 'pi-profile__cel';

  const nowy = document.createElement('button');
  nowy.type = 'button';
  nowy.className = 'dn-btn dn-btn--zarys';
  nowy.textContent = 'Zacznij nowy profil';
  nowy.addEventListener('click', () => wyczysc());

  const element = document.createElement('div');
  element.className = 'pi-profile__formularz';
  element.append(
    cel,
    wiersz('Nazwa profilu', nazwa, { klasa: 'dn-pole' }),
    wiersz('Opis profilu', opis, { klasa: 'dn-pole' }),
    naglowekGrupy('Kontekst — trzy przełączniki (odrębna / współdzielona)'),
    grupaKontekstu,
    naglowekGrupy('Zakres techniczny — osiem zakresów (włączony / wyłączony)'),
    grupaTechniczna,
    nowy,
  );

  function nanieszCel(): void {
    cel.textContent =
      profileId === undefined
        ? 'Zapis założy NOWY profil. Nie przypisuje go do żadnego poziomu — przypisanie to osobna czynność.'
        : `Zapis ZMIENI profil o identyfikatorze ${profileId} we wszystkich jego przypisaniach. ` +
          'Aby zamiast tego założyć nowy, naciśnij „Zacznij nowy profil”.';
  }

  function wyczysc(): void {
    profileId = undefined;
    nazwa.value = '';
    opis.value = '';
    for (const przelacznikPozycji of przelacznikiKontekstu.values()) ustawPrzelacznik(przelacznikPozycji, true, etykietaKontekstu);
    for (const przelacznikPozycji of przelacznikiTechniczne.values()) ustawPrzelacznik(przelacznikPozycji, false, etykietaTechniczna);
    nanieszCel();
  }

  function wypelnij(profil: IsolationProfile): void {
    profileId = profil.id;
    nazwa.value = profil.name;
    opis.value = profil.description ?? '';
    for (const [kind, przelacznikPozycji] of przelacznikiKontekstu) {
      const znaleziony = profil.contextSwitches.find((p) => p.kind === kind);
      ustawPrzelacznik(przelacznikPozycji, znaleziony?.isolated ?? true, etykietaKontekstu);
    }
    for (const [zakres, przelacznikPozycji] of przelacznikiTechniczne) {
      const znaleziony = profil.technicalSwitches.find((p) => p.scope === zakres);
      ustawPrzelacznik(przelacznikPozycji, znaleziony?.isolated ?? false, etykietaTechniczna);
    }
    nanieszCel();
  }

  function zadanie(): IsolationProfileSaveRequest {
    const contextSwitches: IsolationSwitch[] = POZYCJE_KONTEKSTU.map((pozycja) => ({
      kind: pozycja.kind,
      isolated: przelacznikiKontekstu.get(pozycja.kind)?.kontrolka.checked ?? true,
    }));
    const technicalSwitches: IsolationTechnicalSwitch[] = POZYCJE_TECHNICZNE.map((pozycja) => ({
      scope: pozycja.scope,
      isolated: przelacznikiTechniczne.get(pozycja.scope)?.kontrolka.checked ?? false,
    }));

    return {
      ...(profileId === undefined ? {} : { profileId }),
      name: nazwa.value,
      ...(opis.value === '' ? {} : { description: opis.value }),
      contextSwitches,
      technicalSwitches,
    };
  }

  wyczysc();

  return { element, zadanie, wypelnij, wyczysc };
}

function naglowekGrupy(tresc: string): HTMLElement {
  const element = document.createElement('h4');
  element.className = 'pi-profile__naglowek';
  element.textContent = tresc;
  return element;
}

/** Jeden przełącznik: nazwa klucza kontraktu, pole wyboru oraz jego słowna wartość widoczna obok niego. */
function zbudujPrzelacznik(nazwaKlucza: string, etykieta: (isolated: boolean) => string): Przelacznik {
  const kontrolka = przelacznik(nazwaKlucza);

  const wartosc = document.createElement('span');
  wartosc.className = 'dn-plakietka dn-plakietka--informacja';

  const podpis = document.createElement('span');
  podpis.className = 'dn-pole-etykieta';
  podpis.textContent = nazwaKlucza;

  const element = document.createElement('label');
  element.className = 'pi-profile__przelacznik';
  element.append(kontrolka, podpis, wartosc);

  kontrolka.addEventListener('change', () => {
    wartosc.textContent = etykieta(kontrolka.checked);
  });

  return { element, kontrolka };
}

function ustawPrzelacznik(
  przelacznikPozycji: Przelacznik,
  isolated: boolean,
  etykieta: (wartosc: boolean) => string,
): void {
  przelacznikPozycji.kontrolka.checked = isolated;
  przelacznikPozycji.kontrolka.dispatchEvent(new Event('change'));
  const wartosc = przelacznikPozycji.element.querySelector('.dn-plakietka');
  if (wartosc !== null) wartosc.textContent = etykieta(isolated);
}
