import {
  AppComponentKind,
  AppWorkspaceLayer,
  type AppComponent,
  type AppWorkspaceLayer as WarstwaWarsztatu,
} from '../../../../shared/contract';

/**
 * Wykaz komponentów architektury należących do jednej warstwy warsztatu:
 * drzewo komponentów okna Frontend Workspace oraz mapa zależności usług okna
 * Backend Workspace. Wykaz służy tylko do odczytu, a komponenty zmienia się
 * tam, gdzie powstają.
 */
const RODZAJE_WARSTWY: Readonly<Record<WarstwaWarsztatu, readonly string[]>> = {
  [AppWorkspaceLayer.Frontend]: [AppComponentKind.Frontend, AppComponentKind.External],
  [AppWorkspaceLayer.Backend]: [
    AppComponentKind.Backend,
    AppComponentKind.Service,
    AppComponentKind.Database,
    AppComponentKind.Queue,
  ],
};

export interface WykazKomponentowWarstwy {
  element: HTMLElement;
  /** Nanosi komponenty warstwy; oddaje ich liczbę po odsianiu. */
  nanies(komponenty: readonly AppComponent[]): number;
}

export function utworzWykazKomponentowWarstwy(
  warstwa: WarstwaWarsztatu,
  tytul: string,
): WykazKomponentowWarstwy {
  const naglowek = document.createElement('p');
  naglowek.className = 'mp-wykaz__tytul';
  naglowek.textContent = tytul;

  const lista = document.createElement('ul');
  lista.className = 'mp-wykaz';

  const element = document.createElement('div');
  element.className = 'mp-wykaz__powloka';
  element.append(naglowek, lista);

  return {
    element,
    nanies(komponenty) {
      const nasze = komponenty.filter((komponent) =>
        (RODZAJE_WARSTWY[warstwa] ?? []).includes(komponent.kind),
      );
      lista.replaceChildren(...nasze.map(wiersz));
      return nasze.length;
    },
  };
}

/**
 * Jeden komponent warstwy wraz z jego zależnościami. Wiersz niesie nazwę
 * komponentu oraz wyliczenie komponentów, od których zależy, albo adnotację
 * o braku zależności.
 */
function wiersz(komponent: AppComponent): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'mp-wykaz__nazwa';
  nazwa.textContent = komponent.name;

  const zaleznosci = document.createElement('span');
  zaleznosci.className = 'mp-wykaz__zaleznosci';
  const wskazane = komponent.dependsOn ?? [];
  zaleznosci.textContent =
    wskazane.length === 0 ? 'bez zależności' : `zależy od: ${wskazane.join(', ')}`;

  const element = document.createElement('li');
  element.className = 'mp-wykaz__wiersz';
  element.dataset['komponent'] = komponent.id;
  element.append(nazwa, zaleznosci);
  return element;
}

/**
 * Komponenty warstwy odsiane z całego zbioru według rodzaju, przeznaczone do
 * listy wyboru w formularzu zapisu. Kolejność ze zbioru wejściowego zostaje
 * zachowana.
 */
export function komponentyWarstwy(
  warstwa: WarstwaWarsztatu,
  komponenty: readonly AppComponent[],
): readonly AppComponent[] {
  return komponenty.filter((komponent) =>
    (RODZAJE_WARSTWY[warstwa] ?? []).includes(komponent.kind),
  );
}
