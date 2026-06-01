"use client";

import { useCallback, useEffect, useState } from "react";
import { getMyPersonas } from "@/src/utils/api/user/persona";
import { getAccessToken, getActivePersonaID, setActivePersonaID } from "@/src/utils/auth/session";
import type { Persona } from "@/src/types/api/persona";

export function useActivePersona() {
  const [personas, setPersonas] = useState<Persona[]>([]);
  const [activeID, setActiveID] = useState<string>("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getAccessToken();
    if (!token) {
      void Promise.resolve().then(() => setLoading(false));
      return;
    }

    getMyPersonas(token)
      .then((res) => {
        const list = res.data ?? [];
        setPersonas(list);
        const saved = getActivePersonaID();
        const found = list.find((p) => p.id === saved) ?? list[0];
        if (found) {
          setActiveID(found.id);
          setActivePersonaID(found.id);
        }
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const switchPersona = useCallback((id: string) => {
    setActiveID(id);
    setActivePersonaID(id);
  }, []);

  const activePersona = personas.find((p) => p.id === activeID) ?? null;
  return { personas, activeID, activePersona, loading, switchPersona };
}
