import pandas as pd
from sqlalchemy import create_engine, text, Column, Integer, DateTime
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
import logging
from datetime import datetime
from transformation_data import *
from dotenv import load_dotenv
import os

dotenv_path = "./.env"
load_dotenv(dotenv_path)

dw_user = os.getenv('DW_USER')
dw_password = os.getenv('DW_PASS')
dw_host = os.getenv('DW_HOST')
dw_name = os.getenv('DW_NAME')

class PostgreSQLLoader:
    def __init__(self):
        """Inicializa o loader com as credenciais do PostgreSQL."""
        self.connection_string = f'postgresql+pg8000://{dw_user}:{dw_password}@{dw_host}:5433/{dw_name}'
        self.engine = create_engine(self.connection_string)
        self.setup_logging()

    def setup_logging(self):
        """Configura o logging para monitorar o processo de carga."""
        logging.basicConfig(
            level=logging.INFO,
            format='%(asctime)s - %(levelname)s - %(message)s',
            handlers=[
                logging.FileHandler(f'etl_loader_{datetime.now().strftime("%Y%m%d_%H%M%S")}.log'),
                logging.StreamHandler()
            ]
        )
        self.logger = logging.getLogger(__name__)

    def test_connection(self):
        """Testa a conexão com o banco de dados."""
        try:
            with self.engine.connect() as conn:
                self.logger.info("Conexão com o banco de dados estabelecida com sucesso!")
                return True
        except Exception as e:
            self.logger.error(f"Erro ao conectar com o banco de dados: {str(e)}")
            return False

    def get_last_dim_datetime_date(self):
        """Obtém a data do registro mais recente de dim_datetime."""
        query = """
            SELECT date
            FROM dim_datetime
            ORDER BY id DESC
            LIMIT 1;
        """
        try:
            with self.engine.connect() as conn:
                result = conn.execute(text(query))
                last_date = result.scalar()
                
                # Adicionando logs detalhados
                self.logger.info(f"Executando query: {query}")
                self.logger.info(f"Resultado da query: {last_date}")
                
                if last_date is None:
                    self.logger.warning("Nenhuma data encontrada em dim_datetime (primeira execução do ETL).")
                    return None
                    
                self.logger.info(f"Última data obtida em dim_datetime: {last_date}")
                return last_date
        except Exception as e:
            self.logger.error(f"Erro ao obter última data de dim_datetime: {str(e)}")
            raise
        
    def check_updates(self, df, update_column, table_name):
        """
        Verifica atualizações em qualquer tabela comparando a data da coluna de atualização 
        com a última data registrada na tabela dim_datetime.
        
        Args:
            df: DataFrame com os dados.
            update_column: Nome da coluna que contém a data de atualização.
            table_name: Nome da tabela (para logging).
        
        Returns:
            DataFrame contendo apenas os registros atualizados.
        """
        try:
            # Obtém a última data registrada em dim_datetime
            last_dim_datetime_date = self.get_last_dim_datetime_date()

            if update_column and update_column in df.columns:
                if last_dim_datetime_date is not None:
                    # Filtra registros com data maior que a última data registrada
                    updated_records = df[df[update_column] > last_dim_datetime_date]
                    self.logger.info(f"Registros no DataFrame original ({table_name}): {len(df)}")
                    self.logger.info(f"Registros atualizados encontrados ({table_name}): {len(updated_records)}")

                    if updated_records.empty:
                        self.logger.info(f"Nenhum novo registro encontrado para a tabela {table_name}.")
                    return updated_records
                else:
                    self.logger.warning(f"Nenhuma data em dim_datetime. Carregando todos os registros para {table_name}.")
                    return df  # Carrega todos os registros caso seja a primeira execução do ETL
            else:
                self.logger.warning(f"Coluna de atualização '{update_column}' não encontrada em {table_name}. Carregando todos os registros.")
                return df  # Carrega todos os registros caso a coluna não exista ou seja inválida
        except Exception as e:
            self.logger.error(f"Erro ao verificar atualizações para {table_name}: {str(e)}")
            raise


    UPDATE_COLUMNS = {
    'dim_user': 'usr_last_update',
    'dim_process': 'pc_finish_date',
    'dim_vacancy': 'vc_closing_date',
    'hiring_process_candidate' : 'cd_last_update',
    }

    def load_table(self, df, table_name, update_column):
        """Carrega dados em uma tabela após checar atualizações.
        
          Args:
            df: DataFrame com os dados
            table_name: Nome da tabela
            update_column: Coluna de atualização (None para carga completa)
        """
        try:
            # Verifica se é a tabela `dim_datetime`
            if table_name == 'dim_datetime':
                # Se a tabela estiver vazia, insere todos os registros do DataFrame
                with self.engine.connect() as conn:
                    result = conn.execute(text(f"SELECT COUNT(*) FROM {table_name}"))
                    record_count = result.scalar()
                    print(f"Tabela {table_name} contém {record_count} registros antes da carga.")

                if record_count == 0:
                    self.logger.info(f"Tabela {table_name} está vazia. Inserindo todos os registros.")
                    with self.engine.begin() as conn:
                        df.to_sql(
                            name=table_name,
                            con=conn,
                            if_exists='append',
                            index=False,
                            method='multi',
                            chunksize=1000
                        )
                    self.logger.info(f"Carregados {len(df)} registros na tabela {table_name}")
                    return

            updated_df = self.check_updates(df, update_column, table_name)
            if updated_df.empty:
                self.logger.info(f"Não há novos registros para carregar em {table_name}")
                return
            
            with self.engine.begin() as conn:
                updated_df.to_sql(
                    name=table_name,
                    con=conn,
                    if_exists='append',
                    index=False,
                    method='multi',
                    chunksize=1000
                )
                self.logger.info(f"Carregados {len(updated_df)} registros na tabela {table_name}")
        except Exception as e:
            self.logger.error(f"Erro ao carregar tabela {table_name}: {str(e)}")
            raise

    def load_fact_table(self, fact_table):
        """Carrega dados na tabela fato."""
        dim_datetime_id = self.get_dim_datetime_id()
        fact_table['dim_date_id'] = dim_datetime_id
        self.load_table(fact_table, 'fact_hiring_process', update_column=None)

    def main(self):
        """Função principal para teste de inserção."""
        try:
            loader = PostgreSQLLoader()
            if loader.test_connection():
                print("Conexão bem-sucedida! Iniciando carga de teste...")

                # Carrega dim_datetime primeiro, independentemente de outras tabelas
                if dim_datetime is not None:
                    loader.load_table(dim_datetime, 'dim_datetime', update_column=None)

                # Variável para rastrear se alguma tabela foi atualizada
                tables_updated = False

                if dim_department is not None:
                    if loader.load_table(dim_department, 'dim_department', update_column=None):
                        tables_updated = True
                if dim_user is not None:
                    if loader.load_table(dim_user, 'dim_user', 'usr_last_update'):
                        tables_updated = True
                if dim_process is not None:
                    if loader.load_table(dim_process, 'dim_process', 'pc_finish_date'):
                        tables_updated = True
                if dim_vacancy is not None:
                    if loader.load_table(dim_vacancy, 'dim_vacancy', 'vc_closing_date'):
                        tables_updated = True
                if dim_candidate is not None:
                    if loader.load_table(dim_candidate, 'dim_candidate', 'cd_last_update'):
                        tables_updated = True

                # Carrega tabela fato apenas se houve atualização em outras tabelas
                if tables_updated:
                    loader.load_fact_table(fact_hiring_process)
                    print("Carga de teste concluída com sucesso!")
                else:
                    print("Nenhuma atualização detectada, nenhuma carga necessária.")

            else:
                print("--")
        except Exception as e:
            logging.error(f"Erro no processo de teste: {str(e)}")
            raise



if __name__ == "__main__":
    PostgreSQLLoader().main()