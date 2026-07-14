export interface PersonaSetupFields {
  handle: string;
  displayName: string;
  bio: string;
}

export interface PersonaSetupStepProps extends PersonaSetupFields {
  onChange: (fields: PersonaSetupFields) => void;
  onContinue: () => void;
  onBack: () => void;
}

export interface GenreStepProps {
  selected: string[];
  loading: boolean;
  error: string;
  onChange: (tags: string[]) => void;
  onContinue: () => void;
  onBack: () => void;
}
