import { Command, type IsolationProfile } from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import { utworzStanTresci } from '../komponenty/stan-tresci';
import { ETYKIETY_ZASIEGU } from './katalog-izolacji';
import { posijKomende } from './komenda';
import type { ObszarIzolacji, ZaleznosciObszaru } from './obszary';
import { utworzFormularzProfilu } from './profile-formularz';
import { zbudujListeProfili, zbudujPodgladProfilu } from './profile-widok';
import { NAZWY_WARSTW } from './stan-warstwy';

/**
 * Obszar „Profile" — nazwany zestaw trzech przełączników kontekstu i ośmiu
 * zakresów technicznych, z czterema czynnościami rdzenia: Zapisz · Wczytaj ·
 * Przypisz do poziomu · Usuń.
 *
 * Pięć komend rdzenia (`core/kompozycja.go`): `isolation.profile.list` (wykaz),
 * `isolation.profile.save` (zapis nowego albo zmiana istniejącego),
 * `isolation.profile.load` (wczytanie do podglądu), `isolation.profile.assign`
 * (przypisanie do poziomu i warstwy), `isolation.profile.delete` (usunięcie).
 *
 * Wczytanie to nie przypisanie i na tym polega ten obszar. `profile.load`
 * zwraca zawartość profilu do podglądu i do formularza; maszyneria po nim
 * pracuje tak samo jak przedtem. `profile.assign` dopiero wiąże profil
 * z poziomem zasięgu i warstwą i zwraca `IsolationPolicy` — politykę
 * obowiązującą po przypisaniu. Dwie czynności, dwa skutki, dwa osobne zdania
 * przy przyciskach.
 *
 * Cel przypisania jest widoczny wcześniej i nie jest wybierany drugi raz tutaj:
 * poziom zasięgu i identyfikator bytu przychodzą z selektora zasięgu (lewy
 * panel, `stan-zasiegu.ts`), a warstwa z pasa narzędzi okna
 * (`sterowanie-warstwa.ts`) — oba dotyczą wszystkich paneli naraz. Drugi
 * selektor poziomu w tym obszarze pokazywałby cel przypisania inny niż zasięg,
 * na którym Operator właśnie przestawia macierz. Przycisk przy profilu mówi
 * w podpowiedzi, dokąd trafi zapis.
 *
 * Usunięcie idzie od razu, bez „czy na pewno"; żadna pozycja nie jest wygaszona
 * ani zablokowana. Odmowa rdzenia wraca zdaniem trzyczęściowym: co się nie
 * udało, dlaczego (treść wprost z rdzenia) i czym Operator to zmieni.
 */
export function utworzObszar(zaleznosci: ZaleznosciObszaru): ObszarIzolacji {
  const { kanal, warstwa, zasieg } = zaleznosci;
  const tresc = utworzStanTresci('pi');
  const formularz = utworzFormularzProfilu();

  const miejsceListy = document.createElement('div');
  miejsceListy.className = 'pi-profile__miejsce-listy';

  const miejscePodgladu = document.createElement('div');
  miejscePodgladu.className = 'pi-profile__miejsce-podgladu';

  const zapisz = document.createElement('button');
  zapisz.type = 'button';
  zapisz.className = 'dn-btn';
  zapisz.textContent = 'Zapisz profil';
  zapisz.title = 'Zapisuje zestaw z formularza. Nie przypisuje go do żadnego poziomu.';
  zapisz.addEventListener('click', () => void zapiszProfil());

  function opisCelu(): string {
    return `${zasieg.opis()}, warstwa „${NAZWY_WARSTW[warstwa.warstwa()]}”`;
  }

  async function odczytajWykaz(): Promise<void> {
    tresc.ladowanie('Pytam rdzeń o zapisane profile izolacji (isolation.profile.list)…');

    const wynik = await posijKomende(kanal, Command.IsolationProfileList, {});

    const miejsce = tresc.tresc();
    miejsce.append(
      zbudujWstep(),
      zbudujPanelCelu(opisCelu()),
      miejsceListy,
      miejscePodgladu,
      zbudujSekcjeZapisu(),
    );

    if (!wynik.udany || wynik.wynik === undefined) {
      miejsceListy.replaceChildren(
        zdanieOdmowy(
          opisOdmowyBledu('Odczyt wykazu profili', wynik.blad),
          'Wykazu nie widać, ale formularz poniżej działa: zapis nowego profilu jest nadal możliwy, ' +
            'a wykaz odświeży się przy powrocie do tej zakładki.',
        ),
      );
      return;
    }

    if (wynik.wynik.profiles.length === 0) {
      miejsceListy.replaceChildren(
        zdanieZwykle(
          'Rdzeń nie ma jeszcze ani jednego profilu izolacji. To nie jest usterka — profil zakłada ' +
            'się formularzem poniżej i dopiero potem przypisuje do poziomu.',
        ),
      );
      return;
    }

    miejsceListy.replaceChildren(
      zbudujListeProfili(
        wynik.wynik.profiles,
        { wczytaj: (p) => void wczytaj(p), przypisz: (p) => void przypisz(p), usun: (p) => void usun(p) },
        opisCelu(),
      ),
    );
  }

  async function zapiszProfil(): Promise<void> {
    const zadanie = formularz.zadanie();
    const wynik = await posijKomende(kanal, Command.IsolationProfileSave, zadanie);

    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.potwierdzenie(
        `${opisOdmowyBledu('Zapis profilu izolacji', wynik.blad)}. Formularz zachował wpisany zestaw — ` +
          'popraw nazwę albo przełączniki i naciśnij „Zapisz profil” ponownie.',
        false,
      );
      return;
    }

    formularz.wypelnij(wynik.wynik.profile);
    tresc.potwierdzenie(
      `Profil „${wynik.wynik.profile.name}” zapisany. Nie jest jeszcze do niczego przypisany — ` +
        'zrobi to „Przypisz do poziomu” przy jego wierszu.',
      true,
    );
    void odczytajWykaz();
  }

  async function wczytaj(profil: IsolationProfile): Promise<void> {
    const wynik = await posijKomende(kanal, Command.IsolationProfileLoad, { profileId: profil.id });

    if (!wynik.udany || wynik.wynik === undefined) {
      miejscePodgladu.replaceChildren(
        zdanieOdmowy(
          opisOdmowyBledu(`Wczytanie profilu „${profil.name}”`, wynik.blad),
          'Podgląd pozostaje pusty. Naciśnij „Wczytaj” ponownie albo wybierz inny profil — nic nie ' +
            'zostało zmienione w maszynerii, bo wczytanie samo w sobie niczego nie przypisuje.',
        ),
      );
      return;
    }

    formularz.wypelnij(wynik.wynik.profile);
    miejscePodgladu.replaceChildren(zbudujPodgladProfilu(wynik.wynik.profile));
    tresc.potwierdzenie(
      `Profil „${wynik.wynik.profile.name}” wczytany do podglądu i do formularza. ` +
        'Izolacja obowiązująca NIE zmieniła się — wczytanie to nie przypisanie.',
      true,
    );
  }

  async function przypisz(profil: IsolationProfile): Promise<void> {
    const byt = zasieg.bytDoZadania();
    const wynik = await posijKomende(kanal, Command.IsolationProfileAssign, {
      profileId: profil.id,
      scope: zasieg.zasieg(),
      ...(byt === undefined ? {} : { scopeId: byt }),
      layer: warstwa.warstwa(),
    });

    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.potwierdzenie(
        `${opisOdmowyBledu(`Przypisanie profilu „${profil.name}”`, wynik.blad)}. ` +
          'Obowiązuje dalej to, co przedtem. Zmień poziom albo identyfikator bytu w selektorze ' +
          'zasięgu (panel lewy) i naciśnij „Przypisz do poziomu” ponownie.',
        false,
      );
      return;
    }

    const polityka = wynik.wynik.policy;
    tresc.potwierdzenie(
      `Profil „${profil.name}” przypisany: poziom „${ETYKIETY_ZASIEGU[polityka.scope] ?? polityka.scope}”, ` +
        `warstwa „${NAZWY_WARSTW[polityka.layer]}”. Rozstrzygnięcie sprawdzisz w zakładce „Polityka efektywna”.`,
      true,
    );
  }

  async function usun(profil: IsolationProfile): Promise<void> {
    const wynik = await posijKomende(kanal, Command.IsolationProfileDelete, { profileId: profil.id });

    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.potwierdzenie(
        `${opisOdmowyBledu(`Usunięcie profilu „${profil.name}”`, wynik.blad)}. ` +
          'Profil pozostaje na liście. Naciśnij „Usuń” ponownie albo zmień jego zawartość formularzem.',
        false,
      );
      return;
    }

    tresc.potwierdzenie(
      wynik.wynik.deleted
        ? `Profil „${profil.name}” usunięty. Przypisania, które go wskazywały, przestały go wskazywać.`
        : `Rdzeń nie usunął profilu „${profil.name}” i nie podał odmowy — wykaz poniżej pokazuje stan po odczycie.`,
      wynik.wynik.deleted,
    );
    void odczytajWykaz();
  }

  function zbudujPanelCelu(cel: string): HTMLElement {
    const naglowek = document.createElement('h4');
    naglowek.className = 'pi-profile__naglowek';
    naglowek.textContent = 'Cel przypisania';

    const zdanie = document.createElement('p');
    zdanie.className = 'pi-profile__opis';
    zdanie.textContent =
      'Dokąd trafi profil po naciśnięciu „Przypisz do poziomu”. Poziom i byt przychodzą z selektora ' +
      'zasięgu (panel lewy), warstwa z pasa narzędzi okna — oba dotyczą wszystkich paneli naraz, ' +
      'więc ten obszar ich nie wybiera po raz drugi. Poziom globalny nie potrzebuje identyfikatora bytu.';

    const wskazanie = document.createElement('p');
    wskazanie.className = 'pi-profile__cel';
    wskazanie.textContent = `Cel czynny: ${cel}.`;

    const element = document.createElement('div');
    element.className = 'pi-profile__cel-panel';
    element.append(naglowek, zdanie, wskazanie);
    return element;
  }

  function zbudujSekcjeZapisu(): HTMLElement {
    const naglowek = document.createElement('h4');
    naglowek.className = 'pi-profile__naglowek';
    naglowek.textContent = 'Zapis profilu';

    const element = document.createElement('div');
    element.className = 'pi-profile__zapis';
    element.append(naglowek, formularz.element, zapisz);
    return element;
  }

  return {
    element: tresc.element,
    odswiez: () => void odczytajWykaz(),
  };
}

/** Wstęp obszaru: czym jest profil i czym różni się wczytanie od przypisania. */
function zbudujWstep(): HTMLElement {
  const czym = document.createElement('p');
  czym.className = 'pi-profile__opis';
  czym.textContent =
    'Profil to nazwany zestaw jedenastu punktów izolacji: trzech przełączników kontekstu i ośmiu ' +
    'zakresów technicznych. Ten sam profil można przypisać kilku sesjom, rolom, projektom i oknom; ' +
    'zmiana profilu odbija się na wszystkich jego przypisaniach.';

  const roznica = document.createElement('p');
  roznica.className = 'pi-profile__opis';
  roznica.textContent =
    'Wczytanie a przypisanie to dwie różne czynności. „Wczytaj” pokazuje zawartość profilu i wpisuje ' +
    'ją do formularza — maszyneria pracuje dalej tak samo. „Przypisz do poziomu” wiąże profil ' +
    'z wybranym poziomem i warstwą i dopiero wtedy zmienia obowiązującą izolację.';

  const element = document.createElement('div');
  element.className = 'pi-profile__wstep';
  element.append(czym, roznica);
  return element;
}

/** Odmowa w miejscu treści cząstkowej: co się nie udało i czym Operator to zmieni. */
function zdanieOdmowy(co: string, czym: string): HTMLElement {
  const plakietka = document.createElement('span');
  plakietka.className = 'dn-plakietka dn-plakietka--blad';
  plakietka.textContent = 'Nie udało się';

  const zdanie = document.createElement('p');
  zdanie.className = 'pi-profile__opis';
  zdanie.textContent = `${co}. ${czym}`;

  const element = document.createElement('div');
  element.append(plakietka, zdanie);
  return element;
}

function zdanieZwykle(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'pi-profile__opis';
  element.textContent = tresc;
  return element;
}
