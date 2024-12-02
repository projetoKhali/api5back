import pandas as pd
from sqlalchemy import create_engine, text, Column, Integer, DateTime
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
import logging
from datetime import datetime
from transformation_data import *

class PostgreSQLLoader:
    def __init__(self):
        """Inicializa o loader com as credenciais do PostgreSQL."""
        self.connection_string = 'postgresql+pg8000://postgres:pass_warehouse@localhost:5433/postgres'
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

    def get_dim_datetime_id(self):
        """Obtém o id do registro de dim_datetime mais recente."""
        query = """
            SELECT id
            FROM dim_datetime
            ORDER BY id DESC
            LIMIT 1;
        """
        with self.engine.connect() as conn:
            result = conn.execute(text(query))
            dim_datetime_id = result.scalar()
            self.logger.info(f"Obteve id de dim_datetime: {dim_datetime_id}")
            return dim_datetime_id
        

    def get_last_dim_datetime_date(self):
        """Obtém a data do registro mais recente de dim_datetime."""
        query = """
            SELECT date
            FROM dim_datetime
            ORDER BY id DESC
            LIMIT 1;
        """
        with self.engine.connect() as conn:
            result = conn.execute(text(query))
            last_date = result.scalar()
            
            if last_date is None:
                self.logger.warning("Nenhuma data encontrada em dim_datetime (primeira execução do ETL).")
                return None
                
            self.logger.info(f"Obteve última data de dim_datetime: {last_date}")
            return last_date
        
    def check_updates(self, df, update_column, table_name):
        """
        Verifica atualizações em qualquer tabela.
        
        Args:
            df: DataFrame com os dados
            update_column: Nome da coluna que contém a data de atualização
            table_name: Nome da tabela (para logging)
        """
        last_etl_date = self.get_last_dim_datetime_date()
        if last_etl_date is None:
            self.logger.info(f"Primeira execução - carregando todos os registros de {table_name}")
            return df
        
        updated_records = df[df[update_column] > last_etl_date]
        self.logger.info(f"Encontrados {len(updated_records)} registros atualizados em {table_name}")
        return updated_records

    UPDATE_COLUMNS = {
    'dim_user': 'usr_last_update',
    'dim_process': 'pc_finish_date',
    'dim_vacancy': 'vc_closing_date',
    'hiring_process_candidate' : 'cd_last_update',
    }
    def load_table(self, df, table_name, update_column):
        """Carrega dados em uma tabela após checar atualizações."""
        try:
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
                if hiring_process_candidate is not None:
                    if loader.load_table(hiring_process_candidate, 'hiring_process_candidate', 'cd_last_update'):
                        tables_updated = True

                # Carrega dim_datetime apenas se houve atualização em alguma tabela
                if tables_updated:
                    loader.load_table(dim_datetime, 'dim_datetime', update_column=None)
                    loader.load_fact_table(fact_hiring_process)  # Carrega fact_hiring_process em seguida
                    print("Carga de teste concluída com sucesso!")
                else:
                    print("Nenhuma atualização detectada, nenhuma carga necessária.")

            else:
                print("Não foi possível estabelecer conexão com o banco de dados.")
        except Exception as e:
            logging.error(f"Erro no processo de teste: {str(e)}")
            raise


if __name__ == "__main__":
    PostgreSQLLoader().main()